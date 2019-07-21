/** @file
 *****************************************************************************
 * @author     This file is part of czero, developed by sero.cash
 *             and contributors (see AUTHORS).
 * @copyright  MIT license (see LICENSE file)
 *****************************************************************************/


#ifndef LIBZERO_INCLUDE_ZERO_H
#define LIBZERO_INCLUDE_ZERO_H

#ifdef __cplusplus
extern "C" {
#endif

#include "./light.h"
#include "./keys.h"
#include "./output.h"
#include "./license.h"
#include "./info.h"
#include "./input.h"
#include "./balance.h"
#include "./pkg.h"
#include "./input_s.h"

extern void zero_init(const char* account_dir,const unsigned char nettype);

extern void zero_init_inouts();

extern void zero_init_no_circuit();

extern void zero_log_bytes(const unsigned char* bytes,int len);

extern void zero_force_fr(const unsigned char data[32],unsigned char fr[32]);


extern void zero_random32(unsigned char r[32]);

extern void zero_fee_str(char *p);

extern const char* zero_base58_enc(const unsigned char* p,int len);

extern int zero_base58_dec(const char* p,unsigned char* out,int len);


extern void zero_merkle_combine(
    const unsigned char d0[32],
    const unsigned char d1[32],
    unsigned char out[32]
);



extern void zero_hash_0_enter(
    const unsigned char in[40],
    const unsigned char out[64]
);

extern void zero_hash_0_leave(
    const unsigned char in[96],
    const unsigned char out[32]
);

extern void zero_hash_1_enter(
    const unsigned char in[40],
    const unsigned char out[64]
);

extern void zero_hash_1_leave(
    const unsigned char in[96],
    const unsigned char out[32]
);

extern void zero_hash_2_enter(
    const unsigned char in[40],
    const unsigned char out[64]
);

extern void zero_hash_2_leave(
    const unsigned char in[96],
    const unsigned char out[32]
);

#ifdef __cplusplus
}
#endif

#endif //LIBZERO_INCLUDE_ZERO_H
