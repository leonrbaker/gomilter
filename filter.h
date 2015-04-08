/*
Copyright (c) 2015 Leon Baker
This projected is licensed under the terms of the MIT License.
*/

// Set Callback functions in smfiDesc struct
extern void makesmfilter(struct smfiDesc *smfilter);
extern void setConnect(struct smfiDesc *smfilter);
extern void setHelo(struct smfiDesc *smfilter);
extern void setEnvFrom(struct smfiDesc *smfilter);
extern void setEnvRcpt(struct smfiDesc *smfilter);
extern void setHeader(struct smfiDesc *smfilter);
extern void setEoh(struct smfiDesc *smfilter);
extern void setBody(struct smfiDesc *smfilter);
extern void setEom(struct smfiDesc *smfilter);
extern void setAbort(struct smfiDesc *smfilter);
extern void setClose(struct smfiDesc *smfilter);

// Utility functions for things that we can't do in Go

// Return the size of a null terminated pointer array
extern int argv_len(char **argv);

extern int wrap_setmlreply(SMFICTX *ctx, char *rcode, char *xcode, int msgc, char **msgv);
