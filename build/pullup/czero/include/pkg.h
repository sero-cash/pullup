//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_PKG_H
#define LIBCZERO_INCLUDE_PKG_H

#include "constant.h"

extern char zero_pkg(
    //---in---
    const unsigned char key[32],
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char memo[64],
    //---out---
    unsigned char asset_cm_ret[32],
    unsigned char ar_ret[32],
    unsigned char pkg_cm_ret[32],
    unsigned char einfo_ret[ZERO_INFO_WIDTH],
    unsigned char proof_ret[ZERO_PROOF_WIDTH]
);

extern char zero_pkg_verify(
    const unsigned char asset_cm[32],
    const unsigned char pkg_cm[32],
    const unsigned char proof[ZERO_PROOF_WIDTH]
);

extern char zero_pkg_confirm(
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char memo[64],
    const unsigned char ar[32],
    const unsigned char pkg_cm[32]
);

#endif //LIBCZERO_INCLUDE_PKG_H
