import { useEffect, useState } from 'react';
import { ChevronLeft } from 'lucide-react';
import { Link } from 'react-router-dom';
import './AuditPage.css';

interface AuditLog {
  job_id: string;
  tool_id: string;
  exit_code: number;
  created_at: string;
}

export default function AuditPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/admin/audit')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch audit logs');
        return res.json();
      })
      .then(data => {
        setLogs(data || []);
        setLoading(false);
      })
      .catch(err => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <div className="page-container" style={{ padding: '24px' }}>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '24px' }}>
        <Link to="/" style={{ display: 'flex', alignItems: 'center', color: 'var(--color-primary)', textDecoration: 'none', marginRight: '16px' }}>
          <ChevronLeft size={20} /> Back to Catalog
        </Link>
        <h1 className="display" style={{ margin: 0, fontSize: '24px' }}>Audit Log</h1>
      </div>

      {loading && <div>Loading logs...</div>}
      {error && <div className="error-badge">{error}</div>}

      {!loading && !error && (
        <div style={{ overflowX: 'auto', backgroundColor: 'var(--color-surface)', borderRadius: '8px', border: '1px solid var(--color-border)' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--color-border)' }}>
                <th style={{ padding: '12px 16px', color: 'var(--color-text-secondary)', fontWeight: 500 }}>Time</th>
                <th style={{ padding: '12px 16px', color: 'var(--color-text-secondary)', fontWeight: 500 }}>Job ID</th>
                <th style={{ padding: '12px 16px', color: 'var(--color-text-secondary)', fontWeight: 500 }}>Tool ID</th>
                <th style={{ padding: '12px 16px', color: 'var(--color-text-secondary)', fontWeight: 500 }}>Exit Code</th>
              </tr>
            </thead>
            <tbody>
              {logs.length === 0 ? (
                <tr>
                  <td colSpan={4} style={{ padding: '24px', textAlign: 'center', color: 'var(--color-text-secondary)' }}>
                    No audit logs found.
                  </td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr key={log.job_id} style={{ borderBottom: '1px solid var(--color-border)' }}>
                    <td style={{ padding: '12px 16px' }}>{new Date(log.created_at).toLocaleString()}</td>
                    <td style={{ padding: '12px 16px', fontFamily: 'monospace' }}>{log.job_id}</td>
                    <td style={{ padding: '12px 16px' }}>{log.tool_id}</td>
                    <td style={{ padding: '12px 16px' }}>
                      <span style={{ 
                        display: 'inline-block',
                        padding: '2px 8px',
                        borderRadius: '4px',
                        fontSize: '12px',
                        backgroundColor: log.exit_code === 0 ? 'rgba(46, 204, 113, 0.2)' : 'rgba(231, 76, 60, 0.2)',
                        color: log.exit_code === 0 ? '#2ecc71' : '#e74c3c'
                      }}>
                        {log.exit_code === -1 ? 'Error' : log.exit_code}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
