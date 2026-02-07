import { useState } from 'react';
import { NodeType, HttpRequestConfig, TransformConfig, ConditionConfig } from '../../types/workflow';
import Input from '../ui/Input';
import Select from '../ui/Select';
import Button from '../ui/Button';

interface NodeConfigPanelProps {
  nodeId: string;
  nodeType: NodeType;
  config: any;
  onUpdate: (config: any) => void;
  onClose: () => void;
}

export default function NodeConfigPanel({
  nodeId,
  nodeType,
  config,
  onUpdate,
  onClose,
}: NodeConfigPanelProps) {
  const [localConfig, setLocalConfig] = useState(config || {});

  const handleSave = () => {
    onUpdate(localConfig);
    onClose();
  };

  const renderHttpRequestConfig = () => {
    const httpConfig = localConfig as HttpRequestConfig;
    return (
      <div className="space-y-4">
        <Input
          label="URL"
          placeholder="https://api.example.com/endpoint"
          value={httpConfig.url || ''}
          onChange={(e) => setLocalConfig({ ...httpConfig, url: e.target.value })}
        />
        <Select
          label="Method"
          value={httpConfig.method || 'GET'}
          onChange={(e) => setLocalConfig({ ...httpConfig, method: e.target.value as any })}
          options={[
            { value: 'GET', label: 'GET' },
            { value: 'POST', label: 'POST' },
            { value: 'PUT', label: 'PUT' },
            { value: 'DELETE', label: 'DELETE' },
            { value: 'PATCH', label: 'PATCH' },
          ]}
        />
        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700">Headers (JSON)</label>
          <textarea
            className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            rows={4}
            placeholder='{"Content-Type": "application/json"}'
            value={JSON.stringify(httpConfig.headers || {}, null, 2)}
            onChange={(e) => {
              try {
                const headers = JSON.parse(e.target.value);
                setLocalConfig({ ...httpConfig, headers });
              } catch {
                // Invalid JSON, ignore
              }
            }}
          />
        </div>
      </div>
    );
  };

  const renderTransformConfig = () => {
    const transformConfig = localConfig as TransformConfig;
    return (
      <div className="space-y-4">
        <Select
          label="Language"
          value={transformConfig.language || 'javascript'}
          onChange={(e) =>
            setLocalConfig({ ...transformConfig, language: e.target.value as any })
          }
          options={[
            { value: 'javascript', label: 'JavaScript' },
            { value: 'jq', label: 'jq' },
          ]}
        />
        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700">Code</label>
          <textarea
            className="font-mono w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            rows={10}
            placeholder="// Transform the input data&#10;return { ...input, transformed: true };"
            value={transformConfig.code || ''}
            onChange={(e) => setLocalConfig({ ...transformConfig, code: e.target.value })}
          />
        </div>
      </div>
    );
  };

  const renderConditionConfig = () => {
    const conditionConfig = localConfig as ConditionConfig;
    return (
      <div className="space-y-4">
        <Select
          label="Combinator"
          value={conditionConfig.combinator || 'AND'}
          onChange={(e) =>
            setLocalConfig({ ...conditionConfig, combinator: e.target.value as any })
          }
          options={[
            { value: 'AND', label: 'AND' },
            { value: 'OR', label: 'OR' },
          ]}
        />
        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700">Conditions (JSON)</label>
          <textarea
            className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            rows={8}
            placeholder='[{"field": "status", "operator": "equals", "value": "success"}]'
            value={JSON.stringify(conditionConfig.conditions || [], null, 2)}
            onChange={(e) => {
              try {
                const conditions = JSON.parse(e.target.value);
                setLocalConfig({ ...conditionConfig, conditions });
              } catch {
                // Invalid JSON, ignore
              }
            }}
          />
        </div>
      </div>
    );
  };

  return (
    <div className="absolute right-0 top-0 z-50 h-full w-96 border-l border-gray-200 bg-white shadow-xl">
      <div className="flex h-full flex-col">
        <div className="border-b border-gray-200 p-4">
          <h2 className="text-lg font-semibold text-gray-900">Configure Node</h2>
          <p className="text-sm text-gray-500">{nodeType}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-4">
          {nodeType === NodeType.HTTP_REQUEST && renderHttpRequestConfig()}
          {nodeType === NodeType.TRANSFORM && renderTransformConfig()}
          {nodeType === NodeType.CONDITION && renderConditionConfig()}
        </div>

        <div className="border-t border-gray-200 p-4">
          <div className="flex space-x-2">
            <Button onClick={handleSave} className="flex-1">
              Save
            </Button>
            <Button onClick={onClose} variant="outline" className="flex-1">
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
