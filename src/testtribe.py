# -*- coding: utf-8 -*-
"""
Test cases.

Created on Sun Mar 13 21:39:26 2016

@author: John Maloney
"""

import unittest
import tribe
import agent
import rep

class TribeTests(unittest.TestCase):
    
    def test_init(self):
        t = tribe.Tribe(3);
        self.assertEqual(t.total_payouts, 0);
        self.assertEqual(len(t.agents), 3)

    def test_playrounds(self):
        cost = 1
        benefit = 3
        
        t = tribe.Tribe(3)

        # all agents use CO action model
        co = agent.ActionModule([True, False, True, False])
        for a in t.agents:
            a.actm = co
            
        # three rounds will be played
        # total payout will be 12
        self.assertTrue(t.agents[0].chooses_donate(t.agents[1]))
        self.assertTrue(t.agents[0].chooses_donate(t.agents[2]))
        self.assertTrue(t.agents[1].chooses_donate(t.agents[0]))
        self.assertTrue(t.agents[1].chooses_donate(t.agents[2]))
        self.assertTrue(t.agents[2].chooses_donate(t.agents[0]))
        self.assertTrue(t.agents[2].chooses_donate(t.agents[1]))
        self.assertEqual(t.playrounds(cost, benefit), 12)
        
        # test reset
        t.reset()
        self.assertEqual(t.total_payouts, 0);
        for a in t.agents:
            self.assertEqual(a.payout, 0)
        
        # set agent reps to BAD
        for a in t.agents:
            a.rep = rep.BAD

        # three rounds will be played
        # total payout will be 0
        self.assertFalse(t.agents[0].chooses_donate(t.agents[1]))
        self.assertFalse(t.agents[0].chooses_donate(t.agents[2]))
        self.assertFalse(t.agents[1].chooses_donate(t.agents[0]))
        self.assertFalse(t.agents[1].chooses_donate(t.agents[2]))
        self.assertFalse(t.agents[2].chooses_donate(t.agents[0]))
        self.assertFalse(t.agents[2].chooses_donate(t.agents[1]))
        self.assertEqual(t.playrounds(cost, benefit), 6)

    def test_select_parent(self):
        t = tribe.Tribe(3)

        allc = agent.ActionModule([True, True, True, True])
        alld = agent.ActionModule([False, False, False, False])
        
        t.agents[0].payout = -1
        t.agents[0].actm = allc
        t.agents[1].payout = -1
        t.agents[1].actm = allc
        t.agents[2].payout = 10
        t.agents[2].actm = alld
        
        self.assertEqual(t.select_local_parent(), t.agents[2])
        
        t.create_next_gen()
        for a in t.agents:
            self.assertEqual(a.actm, alld)
        
if __name__ == '__main__':
    unittest.main()