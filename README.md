# How to use Decentralized Light Wallet (PC)

[TOC]

## declare

> `Pullup wallet` needs to connect to the full node (gero) above the `v1.0.0-rc2` version, and need to ensure that the full node startup parameters include `--lightNode` and `--rpcapi sero,stake,light`

> **Pullup wallet only supports the purchase of shares voted by StakingNode, and does not support the purchase of SOLO voting shares. **


### Confirm that your system environment has Chrome installed, otherwise the Pullup wallet will not be available.

* Official download page: <https://www.google.com/intl/us-EN/chrome/>



## NOTE

* Relationship between Pullup wallet and PC full-node wallet
  * Do not import accounts between each other as the account files cannot be mixed.
  * Assets can be sent to each other using the deposit addresses.
  * Both can accept withdrawals from exchanges or Flight Wallets. (via deposit address)





## Download

The Pullup wallet is available open source on GitHub and you can download the latest Pullup wallet by visiting the link below

<https://github.com/sero-cash/pullup/releases>





## Installation

* Install for MAC
  * Unzip pullup-mac-xxxtar.gz to pullup file.
  * Copy the pullup to the [Applications] folder.
  * You can launch the Pullup wallet from the launchpad.
  * If blocked by the system, you need to go to [Preferences -> Security and Privacy], click [still open]
* Install for Windows
  * Currently only supports 64-bit win7.1 or higher systems.
  * Unzip pullup-windows-xxxzip to pullup folder
  * Put pullup in any folder
  * Enter the folder and double-click pullup.exe to start the Pullup wallet.
     * Do not change the path of pullup.exe in the folder
  * If the system pops up the firewall interception notification, click [Allow]





## Full Node Selection

There is a card called `Node` on the front page of the wallet.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-f4a44c0339b71fa1.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/400)

Click the two-way arrow on the card to pop up the node selector. You can choose the default node http://129.204.197.105:8545, or you can choose the node you built or built by others, no matter which one you choose. Secure because the wallet will not upload any user privacy information to the node.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-ba2866aa1f6da3fd.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/400)

> If you choose a gero node built by yourself or others, make sure that its gero version number is greater than `v1.0.0-rc2` and that the `--lightNode --rpcapi sero,stake,light` service is enabled.

> The other options are:
>
> * Personal Rpc
> * Set third-party full-node ip and port





## Account Management

### Create a new account

* Pullup wallet uses mnemonic to manage accounts. To create an account, click the button at the bottom of the homepage [Create Account].

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-5eac38ad56db462e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)



* Enter at least 8 characters of any character password twice

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-834d8389b6439235.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* Backup mnemonic

  * A mnemonic is equivalent to your private key. It must be kept safe. Leaking a mnemonic is equivalent to giving the account control to others.

  ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-9ddc22d908a10ed1.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* After entering the home page, you can see the account you just created.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-5bd3b6ecc24fd16a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)



### Backup mnemonic

* The Pullup wallet uses the mnemonic as a medium to back up the account, so you can click on the account box in the home page that needs to be backed up to access the account page.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-ca98c6acc58e7a75.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* After the mnemonic page pops up, please write it down as soon as possible, then close the window and make sure it is safely stored and not leaked

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-e0b664b9fc79c8aa.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)



### Importing mnemonics

* Click the [Use mnemonic import] button below the page for creating an account.

  ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-a4511fbe5196f333.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* The wallet pops up the mnemonic input page, fill in the correct mnemonic on the page, and set the password.

  ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-23ffaff7038a664e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)



## Send and Receive Transactions

### Send transaction

Click the [Send] Tab button to enter the Send Transaction page and fill out the form as required by the page.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-f787ad6927c4e0b6.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* The collection address must be filled in with the 'collect address', and the deposit address is a string of base58 encoded numbers, mainly from the
  * Exchange 'Deposit Address' 
   * Pullup wallet's 'main address' or `collect address' 
   *  Flight Wallet `Deposit Address` 
   *  The`Red`  text line in the full-node wallet account page is  _collection adress_` 
   *  `sero.genPKr`Or  `exchange.getPk`to output address in gero

* After clicking [Send], it will be in the background. `Transaction assembly` has already been signed. Transaction is finally broadcasted to the whole network.

### Receiving Transaction

* The address of the Pullup wallet is divided into the main address and the collection address, both of which can be used as the sending address.
  * ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-638b0e8a5cf32ee3.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)
  * The main address is a special address that can be used multiple times without change.
  * The collection address will change after each use.
  * Only the main address can be used for mining.
  * Both addresses can be used in any other situation:
     * Exchange withdrawals
     * Flight wallet withdrawal
     * Full node wallet transfer
     * Pullup wallet transfer
     * Smart contract transfer

* The biggest problem with the previous full-node wallet was that the transaxtion history could not be displayed and the Pullup wallet can display the transaction history.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-bd8d26c61925b17b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/800)





## Equity Pool (StakingNode)

**SERO's equity pool is a decentralized equity pool that can only be used after 1300000 block height. Registering an equity pool before 1300000 million block height will prompt `Stx Verify Error`. **

The following three PoS mining activities can be carried out in the decentralised light wallet:
 * Equity Pool registration 
 * Purchase of shares voted by equity pool 
 * View my voting income
 For a description of the specific data, please refer to: [How to use gero for Staking] (?file=Tutorial/how-to-staking-using-gero)

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-0a03a9533292bc4f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

### Equity Pool Page

Click [Equity Pool] Tab to enter the equity pool page.

* [My Shares] shows a summary of all share earnings in the account in the wallet

* Click on the [View Details] in the upper right corner to see the details of each share.
* ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-666032bdf1a6ae86.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

*  [Equity Pool] below shows a list of all equity pools and statistics for the entire network.
* Click [Register Equity Pool] to register the equity pool
* Click [Buy Tickets] to authorize an equity pool to vote on behalf of the ticket.
* ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-d0f67df26b37c40d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

### Equity Pool Registration

* You must first register a voting account on the full node.
  * For details on how to register, please refer to [How to use gero for Staking] (?file=Tutorial/how-to-staking-using-gero).
  * The voting account must be permanently unlocked.
* Have 200001 SERO in the Pullup wallet, then click [Register Equity Pool] to enter the registration form.
  * 200000 SERO is used to pledge to create an equity pool, which will be automatically returned after the equity pool service period ends.
  * 1 SERO is used to pay the gas fee for the transaction sent by the registered equity pool. The actual gas cost is much smaller than this value.
* There are a few things to note in the registration form
  * The voting address needs to be filled in with any of the deposit address of the voting account.
    * The address can be obtained from from the command line `sero.genPKr`.
  * The rate must be a number in the range [[25,75]` (minimum draw `25% (i.e 1/4 share of equity pool)`, maximum draw 75%)
  * The account needs to select the account with more than 200,000 SERO.
  * ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-c855b47b3b1acffc.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)
  * Click [Next], enter the password and click [OK]. The wallet will create and send the transaction and will register the equity pool in the chain.
* The registered equity pool will be displayed in the equity pool list after 32 confirmation blocks.

### Buying Shares

It is not enough to have an equity pool (Stake Node). You cannot earn any voting rewards. Users need to buy shares in the equity pool. The equity pool can obtain voting rights at a certain probability when the block is generated. When you get the best equity pool by comparison, you can click on the [Purchase Share] button to buy shares.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-bbcb274b0c6e0df0.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

There are a few things to note when buying:

* Account needs to choose the fund account with assets
* Ideally, users who purchase shares do not need to choose a voting address because the equity pool will vote for them.
  * But if there is a situation where the equity pool is maliciously attacked.
  * At this stage, you can use your own voting account to conduct SOLO voting and reduce your losses.
* The purchase amount needs to be the maximum amount that you can bear and the system will automatically help you buy the share with the current fare.
* Purchase up to 1000 shares per transaction.

![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-2fce062755465e2f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

* Click [Next] to enter the password and then click [OK], you can see the transaction on the account details page.
* ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-bf81536ccb2d659b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)
* You can see the transaction confirmation after about 32 confirmation blocks, and you can see the updated data in the equity pool.
* ![image.png](http://sero-media.s3-website-ap-southeast-1.amazonaws.com/images/jianshu/277023-f85093508acd185b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/600)

## For Developer

- use --rpcHost http://127.0.0.1:8545 set rpc host
- use --webHost http://127.0.0.1:2345 set web host
- use --dev true start dev-mode

