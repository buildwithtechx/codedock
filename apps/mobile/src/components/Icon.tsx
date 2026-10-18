import { MaterialDesignIcons } from '@react-native-vector-icons/material-design-icons/static';
import type React from 'react';
import type { StyleProp, TextStyle } from 'react-native';

export interface IconProps {
  name: string;
  size?: number;
  color?: string;
  style?: StyleProp<TextStyle>;
}

export const Icon: React.FC<IconProps> = ({ name, size = 22, color = '#FAFAFA', style }) => {
  return <MaterialDesignIcons name={name as any} size={size} color={color} style={style} />;
};
