import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Workflow } from '../types/workflow';
import { apiService } from '../services/api';
import { Plus, Play, Trash2 } from 'lucide-react';
import Button from '../components/ui/Button';

export default function WorkflowList() {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadWorkflows();
  }, []);

  const loadWorkflows = async () => {
    try {
      const data = await apiService.getWorkflows();
      setWorkflows(data);
    } catch (error) {
      console.error('Failed to load workflows:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this workflow?')) return;
    try {
      await apiService.deleteWorkflow(id);
      setWorkflows(workflows.filter((w) => w.id !== id));
    } catch (error) {
      console.error('Failed to delete workflow:', error);
    }
  };

  const handleExecute = async (id: string, event: React.MouseEvent) => {
    event.preventDefault();
    try {
      await apiService.executeWorkflow(id);
      alert('Workflow execution started!');
    } catch (error) {
      console.error('Failed to execute workflow:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-gray-500">Loading workflows...</div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-50 p-8">
      <div className="mx-auto max-w-6xl">
        <div className="mb-8 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Workflows</h1>
          <Link to="/workflow/new">
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              New Workflow
            </Button>
          </Link>
        </div>

        {workflows.length === 0 ? (
          <div className="rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
            <h3 className="mb-2 text-lg font-medium text-gray-900">No workflows yet</h3>
            <p className="mb-4 text-gray-500">Get started by creating your first workflow</p>
            <Link to="/workflow/new">
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Create Workflow
              </Button>
            </Link>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {workflows.map((workflow) => (
              <Link
                key={workflow.id}
                to={`/workflow/${workflow.id}`}
                className="block rounded-lg border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md"
              >
                <div className="mb-4">
                  <h3 className="text-lg font-semibold text-gray-900">{workflow.name}</h3>
                  {workflow.description && (
                    <p className="mt-1 text-sm text-gray-500">{workflow.description}</p>
                  )}
                </div>

                <div className="mb-4 text-xs text-gray-500">
                  <div>Nodes: {workflow.nodes.length}</div>
                  <div>Last updated: {new Date(workflow.updatedAt).toLocaleDateString()}</div>
                </div>

                <div className="flex space-x-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={(e) => handleExecute(workflow.id, e)}
                  >
                    <Play className="mr-1 h-3 w-3" />
                    Run
                  </Button>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={(e) => {
                      e.preventDefault();
                      handleDelete(workflow.id);
                    }}
                  >
                    <Trash2 className="h-3 w-3" />
                  </Button>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
