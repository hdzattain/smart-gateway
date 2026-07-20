package milvus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hdzattain/smart-gateway/common"

	milvusclient "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	DefaultAddress = "127.0.0.1:19530"

	SleepMemoryCollection = "user_sleep_memories"
	EmbeddingField        = "embedding"
	UserIDField           = "user_id"
	PayloadField          = "payload"

	EmbeddingDimension = 1536
	PayloadMaxLength   = 8192
)

type Manager struct {
	address string

	mu     sync.RWMutex
	client milvusclient.Client
	closed bool
}

var defaultManager = NewManager(resolveAddress())

func NewManager(address string) *Manager {
	if address == "" {
		address = DefaultAddress
	}
	return &Manager{address: address}
}

func Init(ctx context.Context) error {
	err := defaultManager.Init(ctx)
	go defaultManager.healthLoop()
	return err
}

func GetClient(ctx context.Context) (milvusclient.Client, error) {
	return defaultManager.GetClient(ctx)
}

func Close() error {
	return defaultManager.Close()
}

func (m *Manager) Init(ctx context.Context) error {
	if err := m.connect(ctx); err != nil {
		return err
	}
	return m.ensureSleepMemoryCollection(ctx)
}

func (m *Manager) GetClient(ctx context.Context) (milvusclient.Client, error) {
	m.mu.RLock()
	client := m.client
	closed := m.closed
	m.mu.RUnlock()
	if closed {
		return nil, errors.New("milvus manager is closed")
	}
	if client != nil && m.isHealthy(ctx, client) {
		return client, nil
	}
	if err := m.reconnect(ctx); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.client == nil {
		return nil, errors.New("milvus client is not initialized")
	}
	return m.client, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	if m.client == nil {
		return nil
	}
	err := m.client.Close()
	m.client = nil
	return err
}

func (m *Manager) connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := milvusclient.NewDefaultGrpcClient(ctx, m.address)
	if err != nil {
		return fmt.Errorf("connect milvus at %s: %w", m.address, err)
	}
	if !m.isHealthy(ctx, client) {
		_ = client.Close()
		return fmt.Errorf("milvus at %s is not healthy", m.address)
	}

	m.mu.Lock()
	old := m.client
	m.client = client
	m.closed = false
	m.mu.Unlock()
	if old != nil {
		_ = old.Close()
	}
	return nil
}

func (m *Manager) reconnect(ctx context.Context) error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return errors.New("milvus manager is closed")
	}
	if m.client != nil {
		_ = m.client.Close()
		m.client = nil
	}
	m.mu.Unlock()

	if err := m.connect(ctx); err != nil {
		return err
	}
	return m.ensureSleepMemoryCollection(ctx)
}

func (m *Manager) healthLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.RLock()
		closed := m.closed
		client := m.client
		m.mu.RUnlock()
		if closed {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		healthy := client != nil && m.isHealthy(ctx, client)
		cancel()
		if healthy {
			continue
		}
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		if err := m.reconnect(ctx); err != nil {
			common.SysLog("Milvus reconnect failed: " + err.Error())
		} else {
			common.SysLog("Milvus reconnect succeeded")
		}
		cancel()
	}
}

func (m *Manager) isHealthy(ctx context.Context, client milvusclient.Client) bool {
	if client == nil {
		return false
	}
	state, err := client.CheckHealth(ctx)
	return err == nil && state != nil && state.IsHealthy
}

func (m *Manager) ensureSleepMemoryCollection(ctx context.Context) error {
	client, err := m.currentClient()
	if err != nil {
		return err
	}

	exists, err := client.HasCollection(ctx, SleepMemoryCollection)
	if err != nil {
		return fmt.Errorf("check milvus collection %s: %w", SleepMemoryCollection, err)
	}
	if !exists {
		schema := entity.NewSchema().
			WithName(SleepMemoryCollection).
			WithDescription("Long-term sleep memories isolated by Smart Gateway user ID").
			WithAutoID(true).
			WithField(entity.NewField().
				WithName("id").
				WithDataType(entity.FieldTypeInt64).
				WithIsPrimaryKey(true).
				WithIsAutoID(true)).
			WithField(entity.NewField().
				WithName(UserIDField).
				WithDataType(entity.FieldTypeInt64)).
			WithField(entity.NewField().
				WithName(EmbeddingField).
				WithDataType(entity.FieldTypeFloatVector).
				WithDim(EmbeddingDimension)).
			WithField(entity.NewField().
				WithName(PayloadField).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(PayloadMaxLength))

		if err = client.CreateCollection(ctx, schema, 2); err != nil {
			return fmt.Errorf("create milvus collection %s: %w", SleepMemoryCollection, err)
		}
	}

	if err = m.ensureScalarIndex(ctx, UserIDField); err != nil {
		return err
	}
	if err = m.ensureVectorIndex(ctx); err != nil {
		return err
	}
	if err = client.LoadCollection(ctx, SleepMemoryCollection, false); err != nil {
		return fmt.Errorf("load milvus collection %s: %w", SleepMemoryCollection, err)
	}
	return nil
}

func (m *Manager) ensureScalarIndex(ctx context.Context, field string) error {
	client, err := m.currentClient()
	if err != nil {
		return err
	}
	indexes, err := client.DescribeIndex(ctx, SleepMemoryCollection, field)
	if err == nil && len(indexes) > 0 {
		return nil
	}
	idx := entity.NewGenericIndex("idx_"+field, entity.AUTOINDEX, nil)
	if err = client.CreateIndex(ctx, SleepMemoryCollection, field, idx, false); err != nil {
		return fmt.Errorf("create milvus index on %s.%s: %w", SleepMemoryCollection, field, err)
	}
	return nil
}

func (m *Manager) ensureVectorIndex(ctx context.Context) error {
	client, err := m.currentClient()
	if err != nil {
		return err
	}
	indexes, err := client.DescribeIndex(ctx, SleepMemoryCollection, EmbeddingField)
	if err == nil && len(indexes) > 0 {
		return nil
	}
	idx, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("build milvus vector index config: %w", err)
	}
	if err = client.CreateIndex(ctx, SleepMemoryCollection, EmbeddingField, idx, false); err != nil {
		return fmt.Errorf("create milvus vector index on %s.%s: %w", SleepMemoryCollection, EmbeddingField, err)
	}
	return nil
}

func (m *Manager) currentClient() (milvusclient.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.client == nil {
		return nil, errors.New("milvus client is not initialized")
	}
	return m.client, nil
}

func resolveAddress() string {
	if address := os.Getenv("MILVUS_ADDRESS"); address != "" {
		return address
	}
	return DefaultAddress
}
