//
// Created by tang zhige on 2019/3/13.
//

#ifndef LIBCZERO_INCLUDE_OUTPUT_H
#define LIBCZERO_INCLUDE_OUTPUT_H
#include "constant.h"

extern void zero_out_commitment(
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char memo[ZERO_MEMO_WIDTH],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    const unsigned char rsk[32],
    unsigned char cm[32]
);

extern void zero_root_commitment(
    unsigned long index,
    const unsigned char out_cm[32],
    unsigned char cm[32]
);


extern char zero_output(
    //---in---
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char memo[64],
    const unsigned char pkr[ZERO_PKr_WIDTH],
    int is_v1,
    //---out---
    unsigned char asset_cm_ret[32],
    unsigned char ar_ret[32],
    unsigned char out_cm_ret[32],
    unsigned char einfo_ret[ZERO_INFO_WIDTH],
    unsigned char sbase_ret[32],
    unsigned char proof_ret[ZERO_PROOF_WIDTH]
);

extern void zero_gen_asset_cc(
    //---in---
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    unsigned char asset_cc_ret[32]
);

extern void zero_enc_info(
    //---in---
    const unsigned char key[32],
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char rsk[32],
    const unsigned char memo[64],
    //---out---
    unsigned char einfo_ret[ZERO_INFO_WIDTH]
);

extern char zero_output_verify(
    const unsigned char asset_cm[32],
    const unsigned char out_cm[32],
    const unsigned char rpk[32],
    const unsigned char proof[ZERO_PROOF_WIDTH],
    int is_v1
);

extern char zero_output_confirm(
    const unsigned char tkn_currency[32],
    const unsigned char tkn_value[32],
    const unsigned char tkt_category[32],
    const unsigned char tkt_value[32],
    const unsigned char memo[64],
    const unsigned char pkr[96],
    const unsigned char rsk[64],
    const unsigned char out_cm[32]
);



#endif //LIBCZERO_INCLUDE_OUTPUT_H
