import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:80';

export const api = axios.create({
    baseURL: API_URL,
});

export interface ContainerStatus {
    ip_address: string;
    ping_time: number | null;
    last_success: string | null;
}

export const getContainers = async (): Promise<ContainerStatus[]> => {
    const {data} = await api.get('/container_status');
    console.log(data)
    return data
    // console.log("fetching containers");
    //
    // function getRandomInt(max) {
    //     return Math.floor(Math.random() * max);
    // }
    //
    // let data: any[] = [];
    // var r = getRandomInt(10) +1
    // console.log(r)
    // for (var i = 0; i < r; i++) {
    //     data.push({
    //         ip: `127.0.0.${ i }`,
    //         pingTime: getRandomInt(10_000_000)+100,
    //         lastSuccessful: `2022-11-23T15:55:30.${getRandomInt(1000) +1000}Z`,
    //         status:  getRandomInt(100) > 50 ? 'online' : 'offline'
    //     })
    // }
    // // sleep
    // await new Promise(resolve => setTimeout(resolve, 200));
    //
    // console.log(data)
    // // [
    // //     {
    // //         ip: '127.0.0.1',
    // //         pingTime: 200,
    // //         lastSuccessful: '2022-11-23T15:55:30.000Z',
    // //         status: 'online'
    // //     },
    // //     {
    // //         ip: '127.0.0.1',
    // //         pingTime: 200,
    // //         lastSuccessful: '2022-11-23T15:55:30.000Z',
    // //         status: 'offline'
    // //     }
    // // ];
    // console.log(data)
    // return  data;
};

export const clearContainers = async (): Promise<void> => {
    await api.delete('/container_status');
};