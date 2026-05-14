import type { ToolSpec, Flag, Input } from '../types/models';
import './DynamicForm.css';

interface DynamicFormProps {
  spec: ToolSpec;
  value: Record<string, any>;
  onChange: (value: Record<string, any>) => void;
  errors?: Record<string, string>;
}

export function DynamicForm({ spec, value, onChange, errors = {} }: DynamicFormProps) {
  
  const updateField = (id: string, val: any) => {
    onChange({ ...value, [id]: val });
  };

  const renderFlag = (flag: Flag) => {
    const currentValue = value[flag.id] !== undefined ? value[flag.id] : (flag.default || '');
    
    return (
      <div key={flag.id} className="form-group">
        <label className="form-label label">
          {flag.ui?.label || flag.id}
          {flag.required && <span className="required-star">*</span>}
        </label>
        
        {flag.type === 'boolean' ? (
          <label className="toggle-switch">
            <input 
              type="checkbox" 
              checked={!!currentValue}
              onChange={(e) => updateField(flag.id, e.target.checked)}
            />
            <span className="slider"></span>
          </label>
        ) : flag.type === 'enum' ? (
          <select 
            value={currentValue}
            onChange={(e) => updateField(flag.id, e.target.value)}
          >
            <option value="">Select an option...</option>
            {flag.values?.map(v => <option key={v} value={v}>{v}</option>)}
          </select>
        ) : (
          <input 
            type="text" 
            value={currentValue}
            onChange={(e) => updateField(flag.id, e.target.value)}
            placeholder={`Enter ${flag.ui?.label || flag.id}`}
          />
        )}
        {errors[flag.id] && <div className="field-error" style={{color: '#e74c3c', fontSize: '12px', marginTop: '4px'}}>{errors[flag.id]}</div>}
      </div>
    );
  };

  const renderInput = (input: Input) => {
    const selectedFile = value[input.id] as File | undefined;
    
    return (
      <div key={input.id} className="form-group">
        <label className="form-label label">
          File Input: {input.id}
          {input.required && <span className="required-star">*</span>}
        </label>
        <div className={`file-drop-zone ${selectedFile ? 'has-file' : ''}`}>
          {selectedFile ? (
            <div className="selected-file-info">
              <span className="file-name">{selectedFile.name}</span>
              <button 
                className="clear-file-btn" 
                onClick={(e) => { 
                  e.preventDefault(); 
                  updateField(input.id, undefined); 
                }}
              >
                ×
              </button>
            </div>
          ) : (
            <p className="label">Drag and drop files here, or click to select</p>
          )}
          <input 
            type="file" 
            onChange={(e) => {
              if (e.target.files && e.target.files.length > 0) {
                updateField(input.id, e.target.files[0]);
              }
            }} 
            title={selectedFile ? "Click to change file" : ""}
          />
        </div>
        {errors[input.id] && <div className="field-error" style={{color: '#e74c3c', fontSize: '12px', marginTop: '4px'}}>{errors[input.id]}</div>}
      </div>
    );
  };

  // Group fields by category
  const categories: Record<string, { flags: Flag[], inputs: Input[] }> = {
    'General': { flags: [], inputs: [] }
  };

  spec.flags?.forEach(flag => {
    const cat = flag.ui?.category || 'General';
    if (!categories[cat]) categories[cat] = { flags: [], inputs: [] };
    categories[cat].flags.push(flag);
  });

  spec.inputs?.forEach(input => {
    // Inputs usually don't have ui.category in MVP but let's put them in General or Input
    const cat = 'Input Files';
    if (!categories[cat]) categories[cat] = { flags: [], inputs: [] };
    categories[cat].inputs.push(input);
  });

  return (
    <div className="dynamic-form">
      {Object.entries(categories).map(([catName, items]) => {
        if (items.flags.length === 0 && items.inputs.length === 0) return null;
        
        return (
          <div key={catName} className="form-category">
            <h4 className="category-header label">{catName}</h4>
            <div className="category-content">
              {items.inputs.map(renderInput)}
              {items.flags.map(renderFlag)}
            </div>
          </div>
        );
      })}
    </div>
  );
}
