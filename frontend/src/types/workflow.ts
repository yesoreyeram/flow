import { Node, Edge } from 'reactflow';
import { z } from 'zod';

// Node Type Enums
export enum NodeType {
  HTTP_REQUEST = 'httpRequest',
  TRANSFORM = 'transform',
  CONDITION = 'condition',
  TRIGGER = 'trigger',
  WEBHOOK = 'webhook',
}

// Zod Schemas for validation
export const httpRequestSchema = z.object({
  url: z.string().url(),
  method: z.enum(['GET', 'POST', 'PUT', 'DELETE', 'PATCH']),
  headers: z.record(z.string()).optional(),
  body: z.any().optional(),
  authentication: z
    .object({
      type: z.enum(['none', 'basic', 'bearer', 'apiKey']),
      credentials: z.record(z.string()),
    })
    .optional(),
});

export const transformSchema = z.object({
  code: z.string(),
  language: z.enum(['javascript', 'jq']),
});

export const conditionSchema = z.object({
  conditions: z.array(
    z.object({
      field: z.string(),
      operator: z.enum(['equals', 'notEquals', 'contains', 'greaterThan', 'lessThan']),
      value: z.any(),
    })
  ),
  combinator: z.enum(['AND', 'OR']),
});

export const webhookSchema = z.object({
  path: z.string(),
  method: z.enum(['GET', 'POST', 'PUT', 'DELETE']),
});

// TypeScript Types
export type HttpRequestConfig = z.infer<typeof httpRequestSchema>;
export type TransformConfig = z.infer<typeof transformSchema>;
export type ConditionConfig = z.infer<typeof conditionSchema>;
export type WebhookConfig = z.infer<typeof webhookSchema>;

export type NodeConfig =
  | HttpRequestConfig
  | TransformConfig
  | ConditionConfig
  | WebhookConfig;

export interface WorkflowNode extends Node {
  type: NodeType;
  data: {
    label: string;
    config: NodeConfig;
    outputs?: any;
  };
}

export interface WorkflowEdge extends Edge {}

export interface Workflow {
  id: string;
  name: string;
  description?: string;
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  createdAt: string;
  updatedAt: string;
  version: number;
}

export interface ExecutionResult {
  nodeId: string;
  status: 'success' | 'error' | 'pending';
  output?: any;
  error?: string;
  duration?: number;
}

export interface WorkflowExecution {
  id: string;
  workflowId: string;
  status: 'running' | 'completed' | 'failed';
  startedAt: string;
  completedAt?: string;
  results: ExecutionResult[];
}

export interface ApiError {
  message: string;
  code: string;
  details?: any;
}
