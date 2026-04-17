/** 与 admin records API 的 dataRecordItem 一致 */
export type DataRecordRow = {
  id: number;
  code: string;
  referCode: string;
  source: string;
  route: number;
  routeLabel: string;
  payload: string;
  bizPayload: string;
  createTime?: string;
  updateTime?: string;
};

export type ChannelRow = {
  id: number;
  code: string;
  name: string;
  category: number;
  categoryLabel: string;
  router: string;
  useSecurity: number;
  hasPrivateKey: boolean;
  hasPublicKey: boolean;
  callback: string;
  privateKey?: string;
  publicKey?: string;
};

export type RmqCfg = {
  endpoint?: string;
  appId?: string;
  secret?: string;
  publisher?: { messageGroup?: string; topics?: string[] };
  subscribers?: Array<{ consumerGroup?: string; topic?: string }>;
};
