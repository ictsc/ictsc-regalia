export type ScoreModel = {
  markedScore: number;
  penalty: number;
  score: number;
  maxScore: number;
};

export type ScheduleEntry = {
  name: string;
  startAt: string;
  endAt: string;
};

export type Schedule = {
  hasStarted: boolean;
  current?: ScheduleEntry;
  next?: ScheduleEntry;
  entries: ScheduleEntry[];
};

export type SubmissionStatus = {
  isSubmittable: boolean;
  submittableFrom?: string;
  submittableUntil?: string;
};

export type Problem = {
  code: string;
  title: string;
  maxScore: number;
  category: string;
  sectionSlug?: string;
  score?: ScoreModel;
  submissionableSchedules: ScheduleEntry[];
  submissionStatus?: SubmissionStatus;
};

export type ContestantProfile = {
  name: string;
  displayName: string;
  selfIntroduction: string;
};

export type TeamProfile = {
  code: number;
  name: string;
  organization: string;
  memberLimit: number;
  team: {
    code: number;
    name: string;
    organization: string;
    memberLimit: number;
  };
  members: ContestantProfile[];
};

export type Notice = {
  slug: string;
  body: string;
  title: string;
  markdown: string;
  effectiveFrom: string;
};

export type Rank = {
  rank: number;
  teamCode: number;
  teamName: string;
  organization: string;
  score: number;
  lastEffectiveSubmissionAt?: string;
};

export const DeploymentStatus = {
  UNSPECIFIED: "QUEUED",
  QUEUED: "QUEUED",
  DEPLOYING: "DEPLOYING",
  DEPLOYED: "COMPLETED",
  COMPLETED: "COMPLETED",
  FAILED: "FAILED",
} as const;

export type DeploymentStatus =
  (typeof DeploymentStatus)[keyof typeof DeploymentStatus];
