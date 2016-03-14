# -*- coding: utf-8 -*-
"""
Test cases.

Created on Sun Mar 13 21:16:01 2016

@author: John Maloney
"""

import unittest
import agent
import rep

class AgentTests(unittest.TestCase):
    
    def test_init(self):
        a = agent.Agent();
        self.assertEqual(a.rep, rep.GOOD);
        self.assertEqual(a.payout, 0)
        self.assertEqual(a.num_games, 0)

    def test_playround(self):
        cost = 1
        benefit = 3

        don = agent.Agent();
        self.assertEqual(don.rep, rep.GOOD);
        self.assertEqual(don.payout, 0)
        rec = agent.Agent();
        self.assertEqual(rec.rep, rep.GOOD);
        self.assertEqual(rec.payout, 0)
        
        # configure donor action module
        don.actm = agent.ActionModule([True, False, True, False]);
        
        # GOOD GOOD
        self.assertTrue(don.chooses_donate(rec))
        self.assertEqual(don.playround(rec, cost, benefit), benefit-cost+2*cost);
        self.assertEqual(don.payout, 0)
        self.assertEqual(don.num_games, 1)
        self.assertEqual(rec.payout, 4)
        self.assertEqual(rec.num_games, 1)
        
        # GOOD BAD
        rec.rep = rep.BAD
        self.assertFalse(don.chooses_donate(rec))
        self.assertEqual(don.playround(rec, cost, benefit), 2*cost);
        self.assertEqual(don.payout, 1)
        self.assertEqual(don.num_games, 2)
        self.assertEqual(rec.payout, 5)
        self.assertEqual(rec.num_games, 2)

        # BAD BAD
        don.rep = rep.BAD
        self.assertFalse(don.chooses_donate(rec))
        self.assertEqual(don.playround(rec, cost, benefit), 2*cost);
        self.assertEqual(don.payout, 2)
        self.assertEqual(don.num_games, 3)
        self.assertEqual(rec.payout, 6)
        self.assertEqual(rec.num_games, 3)

        # BAD GOOD
        rec.rep = rep.GOOD
        self.assertTrue(don.chooses_donate(rec))
        self.assertEqual(don.playround(rec, cost, benefit), benefit-cost+2*cost);
        self.assertEqual(don.payout, 2)
        self.assertEqual(don.num_games, 4)
        self.assertEqual(rec.payout, 10)
        self.assertEqual(rec.num_games, 4)

        # reset
        don.reset()
        self.assertEqual(don.payout, 0)
        self.assertEqual(don.num_games, 0)
        rec.reset()
        self.assertEqual(rec.payout, 0)
        self.assertEqual(rec.num_games, 0)
        
if __name__ == '__main__':
    unittest.main()