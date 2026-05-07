import { useEffect, useRef, useState } from 'react';
import { Pin, PinOff } from 'lucide-react';
import './LogStream.css';

interface LogEvent {
  type: string;
  data?: string;
  status?: string;
  exitCode?: number;
}

interface LogStreamProps {
  jobId: string;
  onStatusChange: (status: string) => void;
}

export function LogStream({ jobId, onStatusChange }: LogStreamProps) {
  const [logs, setLogs] = useState<LogEvent[]>([]);
  const [autoScroll, setAutoScroll] = useState(true);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Clear logs on new job
    setLogs([]);
    
    // Connect to SSE
    const eventSource = new EventSource(`/api/jobs/${jobId}/events`);
    
    eventSource.onmessage = (e) => {
      try {
        const event: LogEvent = JSON.parse(e.data);
        setLogs(prev => [...prev, event]);
        
        if (event.type === 'status' && event.status) {
          onStatusChange(event.status);
        } else if (event.type === 'complete') {
          onStatusChange(event.status || (event.exitCode === 0 ? 'succeeded' : 'failed'));
          eventSource.close();
        }
      } catch (err) {
        console.error("Failed to parse SSE message", err);
      }
    };

    eventSource.onerror = () => {
      onStatusChange('failed');
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, [jobId, onStatusChange]);

  useEffect(() => {
    if (autoScroll && containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  return (
    <div className="log-stream-container fade-in">
      <div className="log-header">
        <span className="label" style={{color: 'var(--color-text-primary)'}}>Execution Logs</span>
        <button 
          className="scroll-lock-btn"
          onClick={() => setAutoScroll(!autoScroll)}
          title={autoScroll ? "Disable auto-scroll" : "Enable auto-scroll"}
        >
          {autoScroll ? <Pin size={16} /> : <PinOff size={16} color="var(--color-text-secondary)" />}
        </button>
      </div>
      
      <div className="log-panel" ref={containerRef}>
        {logs.map((log, i) => {
          if (log.type === 'stdout') {
            return <div key={i} className="log-line stdout">{log.data}</div>;
          }
          if (log.type === 'stderr') {
            return (
              <div key={i} className="log-line stderr">
                <span className="log-badge err">ERR</span> {log.data}
              </div>
            );
          }
          if (log.type === 'status' || log.type === 'complete') {
            return (
              <div key={i} className="log-line status">
                <span className="log-badge status">STATUS</span> {log.status || `Exit Code: ${log.exitCode}`}
              </div>
            );
          }
          return null;
        })}
        {logs.length === 0 && <div className="log-line empty">Waiting for logs...</div>}
      </div>
    </div>
  );
}
