import {createRouter,createWebHashHistory} from 'vue-router';
import Index from './pages/Index.vue';
import Team from './pages/Team.vue';
import Lineup from './pages/Lineup.vue';
import BattleInfo from './pages/BattleInfo.vue';

const routes = [
    { 
        path: '/', 
        component: Index ,
        keepalive: true
    },
    {
        path: '/team',
        component: Team
    },
    {
        path: '/battle-info',
        component: BattleInfo
    },
    {
        path: '/lineup',
        component: Lineup
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router;