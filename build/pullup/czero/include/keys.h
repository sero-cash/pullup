//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_KEYS_H
#define LIBCZERO_INCLUDE_KEYS_H

#include "constant.h"

extern void zero_seed2sk(
    const unsigned char seed[32],
    unsigned char sk[ZERO_TK_WIDTH]
);

extern void zero_seed2tk(
    const unsigned char seed[32],
    unsigned char tk[ZERO_TK_WIDTH]
);

extern void zero_seed2pk(
    const unsigned char seed[32],
    unsigned char pk[ZERO_PK_WIDTH]
);


extern void zero_sk2pk(const unsigned char sk[ZERO_PK_WIDTH], unsigned char pk[ZERO_PK_WIDTH]);

extern void zero_sk2tk(const unsigned char sk[ZERO_PK_WIDTH], unsigned char tk[ZERO_TK_WIDTH]);

extern void zero_tk2pk(const unsigned char tk[ZERO_TK_WIDTH], unsigned char pk[ZERO_PK_WIDTH]);

extern void zero_pk2pkr(
    const unsigned char pk[ZERO_PK_WIDTH],
    const unsigned char rnd[32],
    unsigned char pkr[ZERO_PKr_WIDTH]
);

extern char zero_ismy_pkr(
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char tk[ZERO_TK_WIDTH]
);

extern void zero_sign_pkr_by_sk(
    const unsigned char h[32],
    const unsigned char sk[ZERO_PK_WIDTH],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    unsigned char s[64]
);

extern void zero_sign_pkr(
    const unsigned char h[32],
    const unsigned char seed[32],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    unsigned char s[64]
);

extern char zero_verify_pkr(
    const unsigned char h[32],
    const unsigned char s[64],
    const unsigned char pkr[ZERO_PKr_WIDTH]
);


extern void zero_hpkr(
    const unsigned char pkr[ZERO_PKr_WIDTH],
    unsigned char hpkr[ZERO_HPKr_WIDTH]
);

extern char zero_pkr_valid(
    const unsigned char pkr[ZERO_PKr_WIDTH]
);

extern char zero_pk_valid(
    const unsigned char pk[ZERO_PK_WIDTH]
);


#endif //LIBCZERO_INCLUDE_KEYS_H
