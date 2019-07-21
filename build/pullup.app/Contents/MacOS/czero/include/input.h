//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_INPUT_H
#define LIBCZERO_INCLUDE_INPUT_H

#include "constant.h"

extern void zero_til(
    const unsigned char tk[ZERO_TK_WIDTH],
    const unsigned char root_cm[32],
    unsigned char til[32]
);

extern void zero_nil(
    const unsigned char sk[ZERO_TK_WIDTH],
    const unsigned char root_cm[32],
    unsigned char til[32]
);

extern char zero_input(
    //---in---
    const unsigned char seed[32],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char sbase[32],
    const unsigned char einfo[ZERO_INFO_WIDTH],
    unsigned long index,
    const unsigned char anchor[32],
    unsigned long position,
    const unsigned char path[ZERO_PATH_DEPTH*32],
    //---out---
    unsigned char asset_cm_ret[32],
    unsigned char ar_ret[32],
    unsigned char nil_ret[32],
    unsigned char til_ret[32],
    unsigned char proof_ret[ZERO_PROOF_WIDTH]
);

extern char zero_input_by_sk(
    //---in---
    const unsigned char sk[ZERO_PK_WIDTH],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char sbase[32],
    const unsigned char einfo[ZERO_INFO_WIDTH],
    unsigned long index,
    const unsigned char anchor[32],
    unsigned long position,
    const unsigned char path[ZERO_PATH_DEPTH*32],
    //---out---
    unsigned char asset_cm_ret[32],
    unsigned char ar_ret[32],
    unsigned char nil_ret[32],
    unsigned char til_ret[32],
    unsigned char proof_ret[ZERO_PROOF_WIDTH]
);

extern char zero_input_verify(
    const unsigned char asset_cm[32],
    const unsigned char anchor[32],
    const unsigned char nil[32],
    const unsigned char proof[ZERO_PROOF_WIDTH]
);



#endif //LIBCZERO_INCLUDE_INPUT_H
