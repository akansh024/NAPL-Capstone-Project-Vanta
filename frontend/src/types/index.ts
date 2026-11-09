export type Finding={id:number;title:string;description:string;severity:'Critical'|'High'|'Medium'|'Low'|'Info';status:'Open'|'In Progress'|'Resolved';asset:string;evidence:string;remediation:string;source:string;lastSeen:string};
export type Import={id:number;filename:string;scanner:string;status:string;assetCount:number;observationCount:number;importedAt:string};
export type Observation={id:number;asset:string;hostname:string;port:number;protocol:string;service:string;state:string;evidence:string;filename:string;observedAt:string};
export type Dashboard={metrics:Record<string,number>;findings:Finding[];imports:Import[]};
