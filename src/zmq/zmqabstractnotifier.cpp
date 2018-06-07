/*
 Copyright (c) 2009-2010 Satoshi Nakamoto
 Copyright (c) 2009-2017 The Bitcoin Core Developers
 Copyright (c) 2018- The iMorpheus Core Developers
 Distributed under the MIT/X11 software license, see the accompanying
 file COPYING or http://www.opensource.org/licenses/mit-license.php.
*/

#include <zmq/zmqabstractnotifier.h>
#include <util.h>


CZMQAbstractNotifier::~CZMQAbstractNotifier()
{
    assert(!psocket);
}

bool CZMQAbstractNotifier::NotifyBlock(const CBlockIndex * /*CBlockIndex*/)
{
    return true;
}

bool CZMQAbstractNotifier::NotifyTransaction(const CTransaction &/*transaction*/)
{
    return true;
}
