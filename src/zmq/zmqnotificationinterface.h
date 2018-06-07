/*
 Copyright (c) 2009-2010 Satoshi Nakamoto
 Copyright (c) 2009-2017 The Bitcoin Core Developers
 Copyright (c) 2018- The iMorpheus Core Developers
 Distributed under the MIT/X11 software license, see the accompanying
 file COPYING or http://www.opensource.org/licenses/mit-license.php.
*/

#ifndef IMORPHEUS_ZMQ_ZMQNOTIFICATIONINTERFACE_H
#define IMORPHEUS_ZMQ_ZMQNOTIFICATIONINTERFACE_H

#include <validationinterface.h>
#include <string>
#include <map>
#include <list>

class CBlockIndex;
class CZMQAbstractNotifier;

class CZMQNotificationInterface final : public CValidationInterface
{
public:
    virtual ~CZMQNotificationInterface();

    static CZMQNotificationInterface* Create();

protected:
    bool Initialize();
    void Shutdown();

    // CValidationInterface
    void TransactionAddedToMempool(const CTransactionRef& tx) override;
    void BlockConnected(const std::shared_ptr<const CBlock>& pblock, const CBlockIndex* pindexConnected, const std::vector<CTransactionRef>& vtxConflicted, const std::vector<CTransactionRef>& vNameConflicts) override;
    void BlockDisconnected(const std::shared_ptr<const CBlock>& pblock, const CBlockIndex* pindexDelete, const std::vector<CTransactionRef>& vNameConflicts) override;
    void UpdatedBlockTip(const CBlockIndex *pindexNew, const CBlockIndex *pindexFork, bool fInitialDownload) override;

private:
    CZMQNotificationInterface();

    void *pcontext;
    std::list<CZMQAbstractNotifier*> notifiers;
};

#endif // IMORPHEUS_ZMQ_ZMQNOTIFICATIONINTERFACE_H
