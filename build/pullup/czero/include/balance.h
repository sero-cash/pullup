//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_BALANCE_H
#define LIBCZERO_INCLUDE_BALANCE_H

#include "constant.h"

extern void zero_sign_balance(
    //---in---
    int zin_size,
    const unsigned char* zin_acms,
    const unsigned char* zin_ars,
    int zout_size,
    const unsigned char* zout_acms,
    const unsigned char* zout_ars,
    int oin_size,
    const unsigned char* oin_accs,
    int oout_size,
    const unsigned char* oout_accs,
    const unsigned char hash[32],
    //---out---
    unsigned char bsign[64],
    unsigned char bcr[32]
);

extern char zero_verify_balance(
    int zin_size,
    const unsigned char* zin_acms,
    int zout_size,
    const unsigned char* zout_acms,
    int oin_size,
    const unsigned char* oin_accs,
    int oout_size,
    const unsigned char* oout_accs,
    const unsigned char hash[32],
    const unsigned char bcr[32],
    const unsigned char bsign[64]
);

#endif //LIBCZERO_INCLUDE_BALANCE_H
