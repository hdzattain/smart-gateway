/*
Copyright (C) 2025 Smart Gateway

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@smart-gateway.shop
*/

import React, { useEffect, useRef, useState } from 'react';
import { Banner, Button, Col, Form, Row, Spin } from '@douyinfe/semi-ui';
import { Info } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  API,
  removeTrailingSlash,
  showError,
  showSuccess,
} from '../../../helpers';

const toBoolean = (value) => value === true || value === 'true';

export default function SettingsPaymentGatewayMCTPay(props) {
  const { t } = useTranslation();
  const sectionTitle = props.hideSectionTitle ? undefined : t('MCTPay 设置');
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    MCTPayEnabled: false,
    MCTPayMerchantID: '',
    MCTPaySecretKey: '',
    MCTPayCheckoutURL: 'https://mct.com.sg/chn/mctpay/',
    MCTPayWebhookURL: '',
    MCTPayUnitPrice: 1,
    MCTPayMinTopUp: 1,
  });
  const formApiRef = useRef(null);

  useEffect(() => {
    if (props.options && formApiRef.current) {
      const currentInputs = {
        MCTPayEnabled: toBoolean(props.options.MCTPayEnabled),
        MCTPayMerchantID: props.options.MCTPayMerchantID || '',
        MCTPaySecretKey: props.options.MCTPaySecretKey || '',
        MCTPayCheckoutURL:
          props.options.MCTPayCheckoutURL || 'https://mct.com.sg/chn/mctpay/',
        MCTPayWebhookURL: props.options.MCTPayWebhookURL || '',
        MCTPayUnitPrice:
          props.options.MCTPayUnitPrice !== undefined
            ? parseFloat(props.options.MCTPayUnitPrice)
            : 1,
        MCTPayMinTopUp:
          props.options.MCTPayMinTopUp !== undefined
            ? parseFloat(props.options.MCTPayMinTopUp)
            : 1,
      };

      setInputs(currentInputs);
      formApiRef.current.setValues(currentInputs);
    }
  }, [props.options]);

  const handleFormChange = (values) => {
    setInputs(values);
  };

  const submitMCTPay = async () => {
    setLoading(true);
    try {
      const options = [
        { key: 'MCTPayEnabled', value: inputs.MCTPayEnabled },
        { key: 'MCTPayMerchantID', value: inputs.MCTPayMerchantID || '' },
        {
          key: 'MCTPayCheckoutURL',
          value: removeTrailingSlash(inputs.MCTPayCheckoutURL || ''),
        },
        {
          key: 'MCTPayWebhookURL',
          value: removeTrailingSlash(inputs.MCTPayWebhookURL || ''),
        },
        { key: 'MCTPayUnitPrice', value: inputs.MCTPayUnitPrice.toString() },
        { key: 'MCTPayMinTopUp', value: inputs.MCTPayMinTopUp.toString() },
      ];

      if (inputs.MCTPaySecretKey) {
        options.push({ key: 'MCTPaySecretKey', value: inputs.MCTPaySecretKey });
      }

      const requestQueue = options.map((opt) =>
        API.put('/api/option/', {
          key: opt.key,
          value: opt.value,
        }),
      );

      const results = await Promise.all(requestQueue);
      const errorResults = results.filter((res) => !res.data.success);
      if (errorResults.length > 0) {
        errorResults.forEach((res) => {
          showError(res.data.message);
        });
      } else {
        showSuccess(t('更新成功'));
        props.refresh && props.refresh();
      }
    } catch (error) {
      showError(t('更新失败'));
    }
    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      <Form
        initValues={inputs}
        onValueChange={handleFormChange}
        getFormApi={(api) => (formApiRef.current = api)}
      >
        <Form.Section text={sectionTitle}>
          <Banner
            type='info'
            icon={<Info size={16} />}
            description={t(
              'MCTPay 支持 https://mct.com.sg/chn/mctpay/ 收银台。',
            )}
            style={{ marginBottom: 16 }}
          />
          <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Switch field='MCTPayEnabled' label={t('启用 MCTPay')} />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='MCTPayMerchantID'
                label={t('商户 ID')}
                placeholder={t('请输入 MCTPay 商户 ID')}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='MCTPaySecretKey'
                label={t('API 密钥')}
                placeholder={t('敏感信息不会发送到前端显示')}
                type='password'
              />
            </Col>
          </Row>
          <Row
            gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}
            style={{ marginTop: 16 }}
          >
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Input
                field='MCTPayCheckoutURL'
                label={t('收银台地址')}
                placeholder='https://mct.com.sg/chn/mctpay/'
              />
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Input
                field='MCTPayWebhookURL'
                label={t('回调地址')}
                placeholder={t('留空则使用服务器默认回调')}
              />
            </Col>
          </Row>
          <Row
            gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}
            style={{ marginTop: 16 }}
          >
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.InputNumber
                field='MCTPayUnitPrice'
                precision={2}
                label={t('充值价格（本地货币/美金）')}
                placeholder={t('例如：1，就是1本地货币/美金')}
              />
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.InputNumber
                field='MCTPayMinTopUp'
                label={t('最低充值美元数量')}
                placeholder={t('例如：1，就是最低充值1$')}
              />
            </Col>
          </Row>
          <Button onClick={submitMCTPay} style={{ marginTop: 16 }}>
            {t('更新 MCTPay 设置')}
          </Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}
