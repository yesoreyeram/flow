import axios, { AxiosInstance, AxiosError } from 'axios';
import { Workflow, WorkflowExecution, ApiError } from '../types/workflow';

class ApiService {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: '/api',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add auth token if available
        const token = localStorage.getItem('auth_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<ApiError>) => {
        const apiError: ApiError = {
          message: error.response?.data?.message || error.message || 'An error occurred',
          code: error.response?.data?.code || 'UNKNOWN_ERROR',
          details: error.response?.data?.details,
        };
        return Promise.reject(apiError);
      }
    );
  }

  // Workflow APIs
  async getWorkflows(): Promise<Workflow[]> {
    const response = await this.client.get<Workflow[]>('/workflows');
    return response.data;
  }

  async getWorkflow(id: string): Promise<Workflow> {
    const response = await this.client.get<Workflow>(`/workflows/${id}`);
    return response.data;
  }

  async createWorkflow(workflow: Partial<Workflow>): Promise<Workflow> {
    const response = await this.client.post<Workflow>('/workflows', workflow);
    return response.data;
  }

  async updateWorkflow(id: string, workflow: Partial<Workflow>): Promise<Workflow> {
    const response = await this.client.put<Workflow>(`/workflows/${id}`, workflow);
    return response.data;
  }

  async deleteWorkflow(id: string): Promise<void> {
    await this.client.delete(`/workflows/${id}`);
  }

  // Execution APIs
  async executeWorkflow(id: string, input?: any): Promise<WorkflowExecution> {
    const response = await this.client.post<WorkflowExecution>(`/workflows/${id}/execute`, {
      input,
    });
    return response.data;
  }

  async getExecution(executionId: string): Promise<WorkflowExecution> {
    const response = await this.client.get<WorkflowExecution>(`/executions/${executionId}`);
    return response.data;
  }

  async getWorkflowExecutions(workflowId: string): Promise<WorkflowExecution[]> {
    const response = await this.client.get<WorkflowExecution[]>(
      `/workflows/${workflowId}/executions`
    );
    return response.data;
  }

  // Test node execution
  async testNode(nodeType: string, config: any, input?: any): Promise<any> {
    const response = await this.client.post('/nodes/test', {
      nodeType,
      config,
      input,
    });
    return response.data;
  }

  // Health check
  async healthCheck(): Promise<boolean> {
    try {
      await this.client.get('/health');
      return true;
    } catch {
      return false;
    }
  }
}

export const apiService = new ApiService();
