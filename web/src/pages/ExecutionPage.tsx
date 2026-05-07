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
  const [toolSpec] = useState<ToolSpec | null>(null);
  const [formState, setFormState] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Job execution state
  const [activeJobId, setActiveJobId] = useState<string | null>(null);
  const [jobStatus, setJobStatus] = useState<string>('idle'); // idle, running, succeeded, failed, timeout

  useEffect(() => {
    if (!id) return;
    fetch(`/api/admin/tools/drafts`) // For MVP, finding the tool in the list
      .then(res => res.json())
      .then(data => {
        const found = (data || []).find((t: any) => t.ID === id);
        if (found) {
          setRecord(found);
          // Parse YAML content to JSON, or assume API gives JSON
          // We need a proper tools API to return the parsed spec, but for MVP let's assume we can parse it
          // Actually, our API returns Content as string (YAML). 
          // We need to parse it. Let's install js-yaml or just hit a specific tool endpoint.
          // For now, let's just make a mock ToolSpec since we don't have a parse endpoint here
          // OR we can add js-yaml dependency.
        } else {
          setError('Tool not found');
        }
        setLoading(false);
      })
      .catch(err => {
        setError(err.message);
        setLoading(false);
      });
  }, [id]);

  const handleRun = async () => {
    if (!id) return;
    try {
      setJobStatus('running');
      setActiveJobId(null);
      
      const response = await fetch(`/api/jobs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          toolSpecId: id,
          inputs: formState, // Sending the form state
        }),
      });
      
      if (!response.ok) throw new Error('Failed to start job');
      const data = await response.json();
      setActiveJobId(data.id);
    } catch (err: any) {
      setJobStatus('failed');
      alert(`Error starting job: ${err.message}`);
    }
  };

  const handleApplyPreset = (values: Record<string, any>) => {
    setFormState(values);
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
            <button className="btn-secondary">
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
