/*
 Copyright (c) 2009-2010 Satoshi Nakamoto
 Copyright (c) 2009-2017 The Bitcoin Core Developers
 Copyright (c) 2018- The iMorpheus Core Developers
 Distributed under the MIT/X11 software license, see the accompanying
 file COPYING or http://www.opensource.org/licenses/mit-license.php.
*/

#ifndef IMORPHEUS_ZMQ_ZMQPUBLISHNOTIFIER_H
#define IMORPHEUS_ZMQ_ZMQPUBLISHNOTIFIER_H

#include <zmq/zmqabstractnotifier.h>

class CBlockIndex;

class CZMQAbstractPublishNotifier : public CZMQAbstractNotifier
{
private:
    uint32_t nSequence; //!< upcounting per message sequence number

public:

    /* send zmq multipart message
       parts:
          * command
          * data
          * message sequence number
    */
    bool SendMessage(const char *command, const void* data, size_t size);

    bool Initialize(void *pcontext) override;
    void Shutdown() override;
};

class CZMQPublishHashBlockNotifier : public CZMQAbstractPublishNotifier
{
public:
    bool NotifyBlock(const CBlockIndex *pindex) override;
};

class CZMQPublishHashTransactionNotifier : public CZMQAbstractPublishNotifier
{
public:
    bool NotifyTransaction(const CTransaction &transaction) override;
};

class CZMQPublishRawBlockNotifier : public CZMQAbstractPublishNotifier
{
public:
    bool NotifyBlock(const CBlockIndex *pindex) override;
};

class CZMQPublishRawTransactionNotifier : public CZMQAbstractPublishNotifier
{
public:
    bool NotifyTransaction(const CTransaction &transaction) override;
};

#endif // IMORPHEUS_ZMQ_ZMQPUBLISHNOTIFIER_H
