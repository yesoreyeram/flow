import { create } from 'zustand';
import { Node, Edge, applyNodeChanges, applyEdgeChanges, NodeChange, EdgeChange } from 'reactflow';
import { Workflow } from '../types/workflow';

interface WorkflowState {
  // Workflow data
  workflow: Workflow | null;
  nodes: Node[];
  edges: Edge[];

  // Selection
  selectedNode: Node | null;

  // Actions
  setWorkflow: (workflow: Workflow) => void;
  setNodes: (nodes: Node[]) => void;
  setEdges: (edges: Edge[]) => void;
  onNodesChange: (changes: NodeChange[]) => void;
  onEdgesChange: (changes: EdgeChange[]) => void;
  addNode: (node: Node) => void;
  updateNode: (nodeId: string, data: any) => void;
  deleteNode: (nodeId: string) => void;
  selectNode: (node: Node | null) => void;
  addEdge: (edge: Edge) => void;
  reset: () => void;
}

export const useWorkflowStore = create<WorkflowState>((set, get) => ({
  workflow: null,
  nodes: [],
  edges: [],
  selectedNode: null,

  setWorkflow: (workflow) =>
    set({
      workflow,
      nodes: workflow.nodes,
      edges: workflow.edges,
    }),

  setNodes: (nodes) => set({ nodes }),
  setEdges: (edges) => set({ edges }),

  onNodesChange: (changes) => {
    set({
      nodes: applyNodeChanges(changes, get().nodes),
    });
  },

  onEdgesChange: (changes) => {
    set({
      edges: applyEdgeChanges(changes, get().edges),
    });
  },

  addNode: (node) => {
    set({
      nodes: [...get().nodes, node],
    });
  },

  updateNode: (nodeId, data) => {
    set({
      nodes: get().nodes.map((node) =>
        node.id === nodeId ? { ...node, data: { ...node.data, ...data } } : node
      ),
    });
  },

  deleteNode: (nodeId) => {
    set({
      nodes: get().nodes.filter((node) => node.id !== nodeId),
      edges: get().edges.filter((edge) => edge.source !== nodeId && edge.target !== nodeId),
    });
  },

  selectNode: (node) => set({ selectedNode: node }),

  addEdge: (edge) => {
    set({
      edges: [...get().edges, edge],
    });
  },

  reset: () =>
    set({
      workflow: null,
      nodes: [],
      edges: [],
      selectedNode: null,
    }),
}));
