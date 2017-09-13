/*
Copyright (c) 2015 Leon Baker
This projected is licensed under the terms of the MIT License.
*/

#include "libmilter/mfapi.h"
#include "_cgo_export.h"


//#include <stdio.h>   // Testing with fprint remove!!!!


// Set Callback functions in smfiDesc struct

void setConnect(struct smfiDesc *smfilter) {
  smfilter->xxfi_connect = &Go_xxfi_connect;
}

void setHelo(struct smfiDesc *smfilter) {
  smfilter->xxfi_helo = &Go_xxfi_helo;
}

void setEnvFrom(struct smfiDesc *smfilter) {
  smfilter->xxfi_envfrom = &Go_xxfi_envfrom;
}

void setEnvRcpt(struct smfiDesc *smfilter) {
  smfilter->xxfi_envrcpt = &Go_xxfi_envrcpt;
}

void setHeader(struct smfiDesc *smfilter) {
  smfilter->xxfi_header = &Go_xxfi_header;
}

void setData(struct smfiDesc *smfilter) {
  smfilter->xxfi_data = &Go_xxfi_data;
}

void setEoh(struct smfiDesc *smfilter) {
  smfilter->xxfi_eoh = &Go_xxfi_eoh;
}

void setBody(struct smfiDesc *smfilter) {
  smfilter->xxfi_body = &Go_xxfi_body;
}

void setEom(struct smfiDesc *smfilter) {
  smfilter->xxfi_eom = &Go_xxfi_eom;
}

void setAbort(struct smfiDesc *smfilter) {
  smfilter->xxfi_abort = &Go_xxfi_abort;
}

void setClose(struct smfiDesc *smfilter) {
  smfilter->xxfi_close = &Go_xxfi_close;
}

// Utility functions for things that we can't do as easily in Go

// Return the length of a null terminated pointer array
int argv_len(char **argv) {
  int argc = 0;
  while (*argv++ != NULL)
  	++argc;
  return argc;
}

// Wrapper for setmlreply
// Not very elegant way of calling the variadic setmlreply function
int wrap_setmlreply(SMFICTX *ctx, char *rcode, char *xcode, int msgc, char **msgv) {
	switch(msgc) {
		case 1:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], NULL);
		case 2:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], NULL);
		case 3:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				NULL);
		case 4:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], NULL);
		case 5:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], NULL);
		case 6:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], NULL);
		case 7:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], NULL);
		case 8:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], NULL);
		case 9:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], NULL);
		case 10:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				NULL);
		case 11:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], NULL);
		case 12:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], NULL);
		case 13:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], NULL);
		case 14:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], NULL);
		case 15:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], NULL);
		case 16:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], NULL);
		case 17:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				NULL);
		case 18:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], NULL);
		case 19:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], NULL);
		case 20:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], NULL);
		case 21:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], NULL);
		case 22:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], NULL);
		case 23:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], NULL);
		case 24:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				NULL);
		case 25:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], NULL);
		case 26:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], NULL);
		case 27:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], NULL);
		case 28:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], msgv[27], NULL);
		case 29:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], msgv[27], msgv[28], NULL);
		case 30:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], msgv[27], msgv[28], msgv[29], NULL);
		case 31:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], msgv[27], msgv[28], msgv[29], msgv[30],
				NULL);
		case 32:
			return smfi_setmlreply(ctx, rcode, xcode, msgv[0], msgv[1], msgv[2],
				msgv[3], msgv[4], msgv[5], msgv[6], msgv[7], msgv[8], msgv[9],
				msgv[10], msgv[11], msgv[12], msgv[13], msgv[14], msgv[15], msgv[16],
				msgv[17], msgv[18], msgv[19], msgv[20], msgv[21], msgv[22], msgv[23],
				msgv[24], msgv[25], msgv[26], msgv[27], msgv[28], msgv[29], msgv[30],
				msgv[31], NULL);
		default: return -1;
	}

}

 

                
