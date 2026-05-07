import { useState } from 'react';
import type { Preset } from '../types/models';
import './PresetBar.css';

interface PresetBarProps {
  presets: Preset[];
  onSelect: (values: Record<string, any>) => void;
}

export function PresetBar({ presets, onSelect }: PresetBarProps) {
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
        <button 
          key={preset.id}
          className={`preset-pill ${activePreset === preset.id ? 'active' : ''}`}
          onClick={() => handleSelect(preset.id, preset.values)}
        >
          {preset.name}
        </button>
      ))}
    </div>
  );
}
