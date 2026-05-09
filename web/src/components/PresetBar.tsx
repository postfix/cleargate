import { useState } from 'react';
import type { Preset } from '../types/models';
import './PresetBar.css';

interface PresetBarProps {
  presets: (Preset & { isUserDefined?: boolean })[];
  onSelect: (values: Record<string, any>) => void;
  onDelete?: (id: string) => void;
}

export function PresetBar({ presets, onSelect, onDelete }: PresetBarProps) {
  const [activePreset, setActivePreset] = useState<string>('custom');

  const handleSelect = (id: string, values: Record<string, any>) => {
    setActivePreset(id);
    onSelect(values);
  };

  return (
    <div className="preset-bar">
      <button 
        className={`preset-pill ${activePreset === 'custom' ? 'active' : ''}`}
        onClick={() => handleSelect('custom', {})}
      >
        Custom
      </button>
      
      {presets.map(preset => (
        <div key={preset.id} className="preset-pill-container" style={{ display: 'inline-flex', alignItems: 'center' }}>
          <button 
            className={`preset-pill ${activePreset === preset.id ? 'active' : ''}`}
            onClick={() => handleSelect(preset.id, preset.values)}
            style={preset.isUserDefined ? { borderRight: 'none', borderTopRightRadius: 0, borderBottomRightRadius: 0 } : {}}
          >
            {preset.name}
          </button>
          {preset.isUserDefined && onDelete && (
            <button 
              className={`preset-pill delete-btn ${activePreset === preset.id ? 'active' : ''}`}
              onClick={(e) => {
                e.stopPropagation();
                if (window.confirm(`Delete preset "${preset.name}"?`)) {
                  onDelete(preset.id);
                  if (activePreset === preset.id) setActivePreset('custom');
                }
              }}
              style={{ padding: '4px 8px', borderLeft: '1px solid rgba(255,255,255,0.1)', borderTopLeftRadius: 0, borderBottomLeftRadius: 0 }}
              title="Delete preset"
            >
              &times;
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
