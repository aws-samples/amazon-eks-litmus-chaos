import * as React from 'react';
import { PieChart, Legend, Tooltip, Pie, Cell } from 'recharts';

export default function GraphComponent(props) {

  // TODO: Support more colors
  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];
  
  const displayChart = () => {
    
    if (Object.keys(props).length > 0) {
      return (
        <PieChart width={400} height={400}>
            <Legend verticalAlign="top" align="left" />
            <Tooltip />
            <Pie data={props.likes} dataKey="count" nameKey="name" cx="50%" cy="40%" innerRadius={50} outerRadius={80} legendType="circle" label>
              {props.likes.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
          </PieChart>
      )
    } 
  }

  return (
    <>
      {displayChart()}
    </>
  );

  }