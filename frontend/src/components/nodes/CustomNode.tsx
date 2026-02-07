import { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { Globe, Code, GitBranch, Webhook as WebhookIcon } from 'lucide-react';
import { NodeType } from '../../types/workflow';

interface CustomNodeProps {
  data: {
    label: string;
    type: NodeType;
  };
  selected: boolean;
}

const nodeIcons = {
  [NodeType.HTTP_REQUEST]: Globe,
  [NodeType.TRANSFORM]: Code,
  [NodeType.CONDITION]: GitBranch,
  [NodeType.TRIGGER]: WebhookIcon,
  [NodeType.WEBHOOK]: WebhookIcon,
};

const nodeColors = {
  [NodeType.HTTP_REQUEST]: 'bg-blue-100 border-blue-500',
  [NodeType.TRANSFORM]: 'bg-purple-100 border-purple-500',
  [NodeType.CONDITION]: 'bg-yellow-100 border-yellow-500',
  [NodeType.TRIGGER]: 'bg-green-100 border-green-500',
  [NodeType.WEBHOOK]: 'bg-green-100 border-green-500',
};

function CustomNode({ data, selected }: CustomNodeProps) {
  const Icon = nodeIcons[data.type];
  const colorClass = nodeColors[data.type];

  return (
    <div
      className={`min-w-[180px] rounded-lg border-2 p-4 shadow-md ${colorClass} ${
        selected ? 'ring-2 ring-primary-500' : ''
      }`}
    >
      <Handle type="target" position={Position.Top} className="!bg-gray-500" />
      
      <div className="flex items-center space-x-2">
        <Icon className="h-5 w-5" />
        <div>
          <div className="font-semibold text-gray-900">{data.label}</div>
          <div className="text-xs text-gray-500">{data.type}</div>
        </div>
      </div>

      <Handle type="source" position={Position.Bottom} className="!bg-gray-500" />
    </div>
  );
}

export default memo(CustomNode);
