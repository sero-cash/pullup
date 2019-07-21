/** @file
 *****************************************************************************
 * @author     This file is part of czero, developed by sero.cash
 *             and contributors (see AUTHORS).
 * @copyright  MIT license (see LICENSE file)
 *****************************************************************************/

#ifndef LIBZERO_CONSTANT_H
#define LIBZERO_CONSTANT_H

enum {
    ZERO_HPKr_WIDTH=20,
    ZERO_PK_WIDTH=64,
    ZERO_TK_WIDTH=64,
    ZERO_PKr_WIDTH=96,
    ZERO_PATH_DEPTH=29,
    ZERO_PROOF_WIDTH=131,
    ZERO_MEMO_WIDTH=64,
    ZERO_INFO_WIDTH=
            32+ //currency
            32+ //value
            32+ //category
            32+ //value
            32+ //rsk
            64, //memo
    ZERO_LIC_WIDTH=ZERO_PROOF_WIDTH,
};

#endif //LIBZERO_CONSTANT_H
