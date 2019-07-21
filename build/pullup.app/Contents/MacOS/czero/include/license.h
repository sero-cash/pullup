//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_LICENSE_H
#define LIBCZERO_INCLUDE_LICENSE_H

#include "constant.h"

extern char zero_pk2pkr_and_licr(
    //---in---
    const unsigned char pk[ZERO_PK_WIDTH],
    unsigned char pkr[ZERO_PKr_WIDTH],
    unsigned long height,
    //---out---
    unsigned long *counteract_ret,
    unsigned long *limit_l_ret,
    unsigned long *limit_h_ret,
    unsigned char licr_ret[ZERO_LIC_WIDTH]
);

extern char zero_check_licr(
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char licr[ZERO_LIC_WIDTH],
    unsigned long counteract,
    unsigned long limit_l,
    unsigned long limit_h,
    unsigned long height
);


#endif //LIBCZERO_INCLUDE_LICENSE_H
