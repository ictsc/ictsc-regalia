export type Deployment = {
  events: DeploymentEvent[];
  redeployable: boolean;
  maxRedeployment: number;
};

export type DeploymentEvent = {
  revision: number;
  occuredAt: string;
  totalPenalty: number;
  type: string;
  isDeploying: boolean;
};
