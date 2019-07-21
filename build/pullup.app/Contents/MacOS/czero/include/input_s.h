//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_INPUT_S_H
#define LIBCZERO_INCLUDE_INPUT_S_H

#include "constant.h"

extern char zero_input_s(
    //---in---
    const unsigned char ehash[32],
    const unsigned char seed[32],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char root_cm[32],
    //---out---
    unsigned char nil_ret[32],
    unsigned char til_ret[32],
    unsigned char sign_ret[64]
);

extern char zero_input_s_by_sk(
    //---in---
    const unsigned char ehash[32],
    const unsigned char sk[ZERO_PK_WIDTH],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char root_cm[32],
    //---out---
    unsigned char nil_ret[32],
    unsigned char til_ret[32],
    unsigned char sign_ret[64]
);

extern char zero_verify_input_s(
    //---in---
    const unsigned char ehash[32],
    const unsigned char root_cm[32],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char nil[32],
    const unsigned char sign[64]
);


#endif //LIBCZERO_INCLUDE_INPUT_S_H
