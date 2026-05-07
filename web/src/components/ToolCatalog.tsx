import { Link } from 'react-router-dom';
import { Wrench } from 'lucide-react';
import type { ToolSpecRecord } from '../types/models';
import './ToolCatalog.css';

interface ToolCatalogProps {
  tools: ToolSpecRecord[];
}

export function ToolCatalog({ tools }: ToolCatalogProps) {
  return (
    <div className="tool-grid">
      {tools.map((record) => {
        // Parse the YAML/JSON content to get metadata
        // For MVP, assuming backend sends JSON string in Content or we parse it
        // Actually the backend sends Content as YAML string currently because it marshals to YAML
        // In a real app we'd parse it, but for now we have Name and Version from the record.
        return (
          <Link to={`/tool/${record.ID}`} key={record.ID} className="tool-card">
            <div className="tool-card-header">
              <div className="tool-icon">
                <Wrench size={20} color="var(--color-accent)" />
              </div>
              <h3 className="heading">{record.Name}</h3>
            </div>
            <p className="label">Version: {record.Version}</p>
            <div className="tool-tags">
              <span className="tag-pill">{record.Status}</span>
            </div>
          </Link>
        );
      })}
    </div>
  );
}
