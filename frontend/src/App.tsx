import { Routes, Route } from 'react-router-dom';
import WorkflowEditor from './pages/WorkflowEditor';
import WorkflowList from './pages/WorkflowList';
import Layout from './components/Layout';

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<WorkflowList />} />
        <Route path="/workflow/:id" element={<WorkflowEditor />} />
        <Route path="/workflow/new" element={<WorkflowEditor />} />
      </Routes>
    </Layout>
  );
}

export default App;
