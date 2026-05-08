import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { ChevronLeft, Play, Save } from 'lucide-react';
import { DynamicForm } from '../components/DynamicForm';
import { PresetBar } from '../components/PresetBar';
import { LogStream } from '../components/LogStream';
import type { ToolSpec, ToolSpecRecord } from '../types/models';
import './ExecutionPage.css';

export default function ExecutionPage() {
  const { id } = useParams<{ id: string }>();
  const [record, setRecord] = useState<ToolSpecRecord | null>(null);
  const [toolSpec, setToolSpec] = useState<ToolSpec | null>(null);
  const [formState, setFormState] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Job execution state
  const [activeJobId, setActiveJobId] = useState<string | null>(null);
  const [jobStatus, setJobStatus] = useState<string>('idle'); // idle, running, succeeded, failed, timeout

  useEffect(() => {
    if (!id) return;
    fetch('/api/catalog')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch tool');
        return res.json();
      })
      .then(data => {
        const found = (data || []).find((t: any) => t.ID === id);
        if (found) {
          setRecord(found);
          // In a real implementation we would parse found.Content
          // For now, if we don't have a parser, we leave toolSpec null 
          // but we still have the record.
        } else {
          setError('Tool not found');
        }
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to fetch tool:', err);
        setError('Could not load tool details.');
        setLoading(false);
      });
  }, [id]);

  const handleRun = async () => {
    if (!id) return;
    try {
      setJobStatus('running');
      setActiveJobId(null);
      
      // Generate a simple job ID (e.g., using timestamp + random)
      const jobId = `job-${Date.now()}-${Math.floor(Math.random() * 10000)}`;
      setActiveJobId(jobId);

      // Extract files from formState
      const filesToUpload: Record<string, File> = {};
      const nonFileValues: Record<string, any> = {};
      
      for (const [key, value] of Object.entries(formState)) {
        if (value instanceof File) {
          filesToUpload[key] = value;
        } else {
          nonFileValues[key] = value;
        }
      }

      // Upload files if any
      if (Object.keys(filesToUpload).length > 0) {
        const formData = new FormData();
        for (const [key, file] of Object.entries(filesToUpload)) {
          formData.append(key, file);
        }
        
        const uploadRes = await fetch(`/api/upload?job_id=${jobId}`, {
          method: 'POST',
          body: formData,
        });
        
        if (!uploadRes.ok) throw new Error('Failed to upload files');
        
        // Add file names to the values sent to execution
        for (const [key, file] of Object.entries(filesToUpload)) {
          nonFileValues[key] = file.name;
        }
      }

      // Start execution
      const response = await fetch(`/api/execute`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tool_id: id,
          job_id: jobId,
          values: nonFileValues,
        }),
      });
      
      if (!response.ok) throw new Error('Failed to start job');
    } catch (err: any) {
      setJobStatus('failed');
      alert(`Error starting job: ${err.message}`);
    }
  };

  const handleApplyPreset = (values: Record<string, any>) => {
    setFormState(values);
  };

  const handleSavePreset = () => {
    const presetName = prompt("Enter a name for this preset:");
    if (!presetName) return;
    
    // For MVP, we just alert since backend preset API isn't wired
    alert(`[MVP Mode] Saving preset '${presetName}' with values:\n\n${JSON.stringify(formState, null, 2)}\n\nIn production, this would call POST /api/tools/${id}/presets`);
  };

  if (loading) return <div className="page-container">Loading...</div>;
  if (error || !record) return <div className="page-container error-badge">{error || 'Not found'}</div>;

  return (
    <div className="page-container">
      <div className="page-header execution-header">
        <Link to="/" className="back-link">
          <ChevronLeft size={20} /> Back to Catalog
        </Link>
        <div className="tool-title-row">
          <h1 className="display">{record.Name}</h1>
          <span className="tag-pill">{record.Version}</span>
        </div>
      </div>

      {toolSpec && toolSpec.presets && toolSpec.presets.length > 0 && (
        <PresetBar presets={toolSpec.presets} onSelect={handleApplyPreset} />
      )}

      <div className="execution-content">
        <div className="form-section">
          {toolSpec ? (
            <DynamicForm 
              spec={toolSpec} 
              value={formState} 
              onChange={setFormState} 
            />
          ) : (
            <div className="form-placeholder">
              <p className="label">Form dynamically generated from ToolSpec</p>
              {/* Fallback rendering since we can't parse YAML purely in browser easily without js-yaml yet */}
              <p className="code" style={{marginTop: '10px'}}>{record.Content.substring(0, 200)}...</p>
            </div>
          )}

          <div className="form-actions">
            <button 
              className="btn-primary" 
              onClick={handleRun}
              disabled={jobStatus === 'running'}
            >
              <Play size={16} style={{marginRight: '8px', verticalAlign: 'middle'}}/>
              Run Tool
            </button>
            <button className="btn-secondary" onClick={handleSavePreset}>
              <Save size={16} style={{marginRight: '8px', verticalAlign: 'middle'}}/>
              Save as Preset
            </button>
          </div>
        </div>
      </div>

      {/* Status Bar */}
      <div className="status-bar">
        <div className="status-indicator">
          <span className={`status-dot ${jobStatus}`}></span>
          <span className="label" style={{textTransform: 'uppercase'}}>{jobStatus}</span>
        </div>
        {activeJobId && <span className="code">Job ID: {activeJobId}</span>}
      </div>

      {/* Log Stream */}
      {activeJobId && (
        <LogStream jobId={activeJobId} onStatusChange={setJobStatus} />
      )}
    </div>
  );
}
