import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';

import { sequence } from '0xsequence'
import feet from './feet.png'
import trunks from './trunks.png'
import body from './body.png'
import head from './head.png'

let items = [
  head,
  body,
  trunks,
  feet
]

const Selected = (props: any) => {
  const squareStyle = {
    border: "1px solid black",
    width: "50px",
    height: "50px",
  };

  const renderGrid = (items: any) => {
    const columnItems = items.map((src: any, index: number) => (
      <div key={index} style={squareStyle}>
        <img width={'47px'} style={{ padding: '1px' }} src={src} alt="Item" />
      </div>
    ));

    return (
      <div style={{ display: "flex", flexDirection: "column" }}>
        {columnItems}
      </div>
    );
  };

  return <div className="grid-container">{renderGrid(props.items)}</div>;
};

const Grid = (props: any) => {
  const gridRows = 3;
  const gridColumns = 9;

  const equip = (src: any, slot: number) => {
    const tempArray = props.shell.slice(); // Create a copy of the shell array
    tempArray[slot] = src;
    props.setShell(tempArray);
  };

  const handleSquareClick = (row: any, col: any) => {
    const index = row * gridColumns + col;
    return index; // console.log(index);
  };

  const renderGrid = () => {
    const grid = [];
    const items_ = [
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1122961902977437778/morganmoskalyk_a_gold_laced_rock_on_a_white_background_9a7ab133-418b-4a7e-83ba-e2db334f4020.png?width=800&height=800', slot: 1 },
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1121170416501784686/morganmoskalyk_a_crystal_talisman_with_ornate_solarpunk_with_th_ebee3978-9530-4dcd-8426-10d40504872a.png?width=800&height=800', slot: 1 },
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1121170310918586378/morganmoskalyk_running_shoes_with_solarpunk_style_and_bubbly_bo_1d386767-4dbc-4168-92f7-d966e9407151.png?width=800&height=800', slot: 3 },
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1122983118702391326/morganmoskalyk_blue_jeans_with_tears_on_a_white_background_41b1e023-5b5c-4486-ac75-03cd78922230.png?width=800&height=800', slot: 2},
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1122985412621783170/morganmoskalyk_crystal_and_metallic_sunglasses_on_a_white_backg_66f96373-2d7d-453c-906d-277db608379c.png?width=800&height=800', slot: 0},
      { srcA: 'https://media.discordapp.net/attachments/1091126083895693312/1121170973048180777/morganmoskalyk_a_male_brimmed_hat_with_this_style_f9374958-cba6-4dff-95b9-a61ab55c985d.png?width=800&height=800', slot: 0}
    ];

    let i = 0;
    for (let row = 0; row < gridRows; row++) {
      const rowItems: any = [];

      for (let col = 0; col < gridColumns; col++) {
        const squareStyle = {
          border: "1px solid black",
          width: "50px",
          height: "50px",
          cursor: 'pointer'
        };
        rowItems.push(
          <div key={col} style={squareStyle} >
            {
              items_[i] ? <img
                width={'50px'}
                style={{ padding: '1px' }}
                src={items_[i]!.srcA!}
                alt="Item"
                onClick={() => {
                  equip(items_[handleSquareClick(row, col)]!.srcA!, items_[handleSquareClick(row, col)]!.slot!)
                }}
              />
                : null
            }
          </div>);
        i++;
      }

      grid.push(
        <div key={row} style={{ display: "flex" }}>
          {rowItems}
        </div>
      );
    }

    return grid;
  };

  return <div className="grid-container">{renderGrid()}</div>;
};


function App() {
  const [shell, setShell] = useState<any>([
    head,
    body,
    trunks,
    feet
  ]);

  sequence.initWallet('mumbai');

  const login = async () => {
    const wallet = sequence.getWallet();
    const connectWallet = await wallet.connect({
      networkId: 80001,
      app: 'test',
      authorize: true,
      settings: {
        theme: 'dark'
      }
    });

    console.log(connectWallet);
  };

  const construct = async () => {
    const wallet = sequence.getWallet();
    const signer = wallet.getSigner()
    const message = 'Include'

    const signature = await signer.signMessage(message)
    console.log(signature)
  }

  React.useEffect(() => {
    console.log('testing');
  }, [shell]);

  return (
    <div className="App">
      <button onClick={() => login()}>login</button>
      <button onClick={() => construct()}>construct</button>
      <div className="flex-container">
        <div className="flex-item"><Selected items={shell} /><br/><Grid shell={shell} setShell={setShell} /></div>
      </div>
      <br />
      <br />
      <br />
    </div>
  );
}

export default App;
