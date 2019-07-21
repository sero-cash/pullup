//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_INFO_H
#define LIBCZERO_INCLUDE_INFO_H
#include "constant.h"

extern char zero_fetch_key(
    const unsigned char tk[ZERO_TK_WIDTH],
    const unsigned char rpk[32],
    const unsigned char key[32]
);

extern void zero_dec_einfo(
    //---in---
    const unsigned char key[32],
    char flag,
    const unsigned char einfo[ZERO_INFO_WIDTH],
    //---out---
    unsigned char tkn_currency_ret[32],
    unsigned char tkn_value_ret[32],
    unsigned char tkt_category_ret[32],
    unsigned char tkt_value_ret[32],
    unsigned char rsk_ret[32],
    unsigned char memo_ret[64]
);

#endif //LIBCZERO_INCLUDE_INFO_H
