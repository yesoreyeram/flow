import { useCallback, useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ReactFlow, {
  Controls,
  Background,
  MiniMap,
  Connection,
  addEdge,
  BackgroundVariant,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { useWorkflowStore } from '../stores/workflowStore';
import { NodeType } from '../types/workflow';
import CustomNode from '../components/nodes/CustomNode';
import NodeConfigPanel from '../components/editor/NodeConfigPanel';
import { Plus, Play, Save } from 'lucide-react';
import Button from '../components/ui/Button';
import { apiService } from '../services/api';

const nodeTypes = {
  [NodeType.HTTP_REQUEST]: CustomNode,
  [NodeType.TRANSFORM]: CustomNode,
  [NodeType.CONDITION]: CustomNode,
  [NodeType.TRIGGER]: CustomNode,
  [NodeType.WEBHOOK]: CustomNode,
};

export default function WorkflowEditor() {
  const { id } = useParams();
  const {
    nodes,
    edges,
    selectedNode,
    setWorkflow,
    onNodesChange,
    onEdgesChange,
    addNode,
    addEdge: addEdgeToStore,
    selectNode,
    updateNode,
  } = useWorkflowStore();

  const [showNodeMenu, setShowNodeMenu] = useState(false);
  const [showConfig, setShowConfig] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (id && id !== 'new') {
      loadWorkflow(id);
    }
  }, [id]);

  const loadWorkflow = async (workflowId: string) => {
    try {
      const workflow = await apiService.getWorkflow(workflowId);
      setWorkflow(workflow);
    } catch (error) {
      console.error('Failed to load workflow:', error);
    }
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const workflow = {
        name: 'My Workflow',
        nodes,
        edges,
      };

      if (id && id !== 'new') {
        await apiService.updateWorkflow(id, workflow);
      } else {
        await apiService.createWorkflow(workflow);
      }
    } catch (error) {
      console.error('Failed to save workflow:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleExecute = async () => {
    if (!id || id === 'new') return;
    try {
      await apiService.executeWorkflow(id);
      alert('Workflow execution started!');
    } catch (error) {
      console.error('Failed to execute workflow:', error);
      alert('Failed to execute workflow');
    }
  };

  const onConnect = useCallback(
    (connection: Connection) => {
      const edge = {
        ...connection,
        id: `e${connection.source}-${connection.target}`,
        type: 'default',
      };
      addEdgeToStore(edge as any);
    },
    [addEdgeToStore]
  );

  const handleAddNode = (type: NodeType) => {
    const newNode = {
      id: `node-${Date.now()}`,
      type,
      position: { x: 250, y: 250 },
      data: {
        label: `${type} Node`,
        type,
        config: {},
      },
    };
    addNode(newNode);
    setShowNodeMenu(false);
  };

  const handleNodeClick = (_event: any, node: any) => {
    selectNode(node);
    setShowConfig(true);
  };

  const handleConfigUpdate = (config: any) => {
    if (selectedNode) {
      updateNode(selectedNode.id, { config });
    }
  };

  return (
    <div className="relative h-full w-full">
      {/* Toolbar */}
      <div className="absolute left-4 top-4 z-10 flex space-x-2">
        <Button onClick={() => setShowNodeMenu(!showNodeMenu)} size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Node
        </Button>
        <Button onClick={handleSave} disabled={isSaving} size="sm" variant="secondary">
          <Save className="mr-2 h-4 w-4" />
          {isSaving ? 'Saving...' : 'Save'}
        </Button>
        <Button onClick={handleExecute} size="sm" variant="outline">
          <Play className="mr-2 h-4 w-4" />
          Execute
        </Button>
      </div>

      {/* Node Menu */}
      {showNodeMenu && (
        <div className="absolute left-4 top-20 z-10 w-64 rounded-lg border border-gray-200 bg-white shadow-lg">
          <div className="p-2">
            <h3 className="mb-2 px-2 text-sm font-semibold text-gray-700">Add Node</h3>
            <div className="space-y-1">
              {Object.values(NodeType).map((type) => (
                <button
                  key={type}
                  onClick={() => handleAddNode(type)}
                  className="w-full rounded px-3 py-2 text-left text-sm hover:bg-gray-100"
                >
                  {type}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* React Flow Canvas */}
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={handleNodeClick}
        nodeTypes={nodeTypes}
        fitView
      >
        <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
        <Controls />
        <MiniMap />
      </ReactFlow>

      {/* Node Configuration Panel */}
      {showConfig && selectedNode && (
        <NodeConfigPanel
          nodeId={selectedNode.id}
          nodeType={selectedNode.data.type}
          config={selectedNode.data.config}
          onUpdate={handleConfigUpdate}
          onClose={() => setShowConfig(false)}
        />
      )}
    </div>
  );
}
