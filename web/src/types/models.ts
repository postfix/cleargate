export interface ToolSpec {
  metadata: {
    name: string;
    displayName: string;
    description: string;
    version: string;
    tags?: string[];
  };
  flags?: Flag[];
  inputs?: Input[];
  presets?: Preset[];
}

export interface Flag {
  id: string;
  type: string;
  required?: boolean;
  default?: any;
  values?: string[];
  ui?: {
    label: string;
    category?: string;
    widget?: string;
  };
}

export interface Input {
  id: string;
  type: string;
  required: boolean;
  destination: string;
  maxSizeMB?: number;
  allowedExtensions?: string[];
}

export interface Preset {
  id: string;
  tool_id?: string;
  name: string;
  values: Record<string, any>;
}

export interface ToolSpecRecord {
  ID: string;
  Name: string;
  Version: string;
  Status: string;
  Content: string;
  CreatedAt: string;
}
