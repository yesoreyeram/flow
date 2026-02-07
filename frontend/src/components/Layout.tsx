import { ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { Workflow, Plus } from 'lucide-react';

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex h-screen flex-col">
      {/* Header */}
      <header className="border-b border-gray-200 bg-white">
        <div className="flex h-16 items-center justify-between px-6">
          <Link to="/" className="flex items-center space-x-2">
            <Workflow className="h-8 w-8 text-primary-600" />
            <span className="text-xl font-bold text-gray-900">Flow</span>
          </Link>
          <nav className="flex items-center space-x-4">
            <Link
              to="/workflow/new"
              className="flex items-center space-x-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
            >
              <Plus className="h-4 w-4" />
              <span>New Workflow</span>
            </Link>
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-hidden">{children}</main>
    </div>
  );
}
