import { useEffect, useState } from 'react';
import { ToolCatalog } from '../components/ToolCatalog';
import type { ToolSpecRecord } from '../types/models';

export default function CatalogPage() {
  const [tools, setTools] = useState<ToolSpecRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/catalog')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch tools');
        return res.json();
      })
      .then(data => {
        setTools(data || []);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to load catalog:', err);
        setError('Could not connect to backend to load tools.');
        setLoading(false);
      });
  }, []);

  return (
    <div className="page-container">
      <div className="page-header">
        <h1 className="display">Tool Catalog</h1>
        <p className="label">Select an approved tool to execute.</p>
      </div>
      
      {loading && <p>Loading tools...</p>}
      {error && <div className="error-badge">Error: {error}</div>}
      
      {!loading && !error && tools.length === 0 && (
        <div className="empty-state">
          <h2 className="heading">No tools available</h2>
          <p>No approved tools found. Contact your administrator to add tools.</p>
        </div>
      )}

      {!loading && tools.length > 0 && <ToolCatalog tools={tools} />}
    </div>
  );
}
