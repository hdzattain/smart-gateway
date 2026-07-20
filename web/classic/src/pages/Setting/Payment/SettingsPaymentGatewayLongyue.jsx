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

export default function SettingsPaymentGatewayLongyue(props) {
  const { t } = useTranslation();
  const sectionTitle = props.hideSectionTitle ? undefined : t('龙跃外卡设置');
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    LongyueEnabled: false,
    LongyueAppId: '',
    LongyueSecretKey: '',
    LongyueApiBase: '',
    LongyueUnitPrice: 1,
    LongyueMinTopUp: 1,
    LongyueCurrency: 'USD',
  });
  const formApiRef = useRef(null);

  useEffect(() => {
    if (props.options && formApiRef.current) {
      const currentInputs = {
        LongyueEnabled: toBoolean(props.options.LongyueEnabled),
        LongyueAppId: props.options.LongyueAppId || '',
        LongyueSecretKey: props.options.LongyueSecretKey || '',
        LongyueApiBase: props.options.LongyueApiBase || '',
        LongyueUnitPrice:
          props.options.LongyueUnitPrice !== undefined
            ? parseFloat(props.options.LongyueUnitPrice)
            : 1,
        LongyueMinTopUp:
          props.options.LongyueMinTopUp !== undefined
            ? parseFloat(props.options.LongyueMinTopUp)
            : 1,
        LongyueCurrency: props.options.LongyueCurrency || 'USD',
      };

      setInputs(currentInputs);
      formApiRef.current.setValues(currentInputs);
    }
  }, [props.options]);

  const handleFormChange = (values) => {
    setInputs(values);
  };

  const submitLongyue = async () => {
    setLoading(true);
    try {
      const options = [
        { key: 'LongyueEnabled', value: inputs.LongyueEnabled },
        { key: 'LongyueAppId', value: inputs.LongyueAppId || '' },
        {
          key: 'LongyueApiBase',
          value: removeTrailingSlash(inputs.LongyueApiBase || ''),
        },
        { key: 'LongyueUnitPrice', value: inputs.LongyueUnitPrice.toString() },
        { key: 'LongyueMinTopUp', value: inputs.LongyueMinTopUp.toString() },
        { key: 'LongyueCurrency', value: inputs.LongyueCurrency || 'USD' },
      ];

      if (inputs.LongyueSecretKey) {
        options.push({ key: 'LongyueSecretKey', value: inputs.LongyueSecretKey });
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
              '龙跃外卡收单系统（CNP），支持国际信用卡支付，用户将跳转到龙跃收银台完成支付。',
            )}
            style={{ marginBottom: 16 }}
          />
          <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Switch field='LongyueEnabled' label={t('启用龙跃外卡支付')} />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='LongyueAppId'
                label={t('商户 AppID')}
                placeholder={t('请输入龙跃商户 AppID')}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='LongyueSecretKey'
                label={t('签名密钥')}
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
                field='LongyueApiBase'
                label={t('API 基础地址')}
                placeholder={t('例如：https://pay.example.com')}
              />
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Input
                field='LongyueCurrency'
                label={t('结算货币')}
                placeholder={t('例如：USD')}
              />
            </Col>
          </Row>
          <Row
            gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}
            style={{ marginTop: 16 }}
          >
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.InputNumber
                field='LongyueUnitPrice'
                precision={2}
                label={t('充值价格（结算货币/美金）')}
                placeholder={t('例如：1，就是1结算货币/美金')}
              />
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.InputNumber
                field='LongyueMinTopUp'
                label={t('最低充值美元数量')}
                placeholder={t('例如：1，就是最低充值1$')}
              />
            </Col>
          </Row>
          <Button onClick={submitLongyue} style={{ marginTop: 16 }}>
            {t('更新龙跃外卡设置')}
          </Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}
