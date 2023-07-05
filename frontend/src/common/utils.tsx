import { red, volcano, gold, yellow, lime, green, cyan, blue, geekblue, purple, magenta, grey } from '@ant-design/colors';

export const hashCode = (str: string) => {
    // java String#hashCode
    var hash = 0;
    for (var i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    return hash;
};

export const intToColor = (i: number) => {
    const first = 3
    const last = 7
    //  const colors = ['#f56a00', '#7265e6', '#ffbf00', '#00a2ae', '#87d068'];
    // const colors = [...red, ...volcano, ...gold, ...yellow, ...lime, ...green, ...cyan, ...blue, ...geekblue, ...purple, ...magenta, ...grey];
    const colors = [...red.slice(first, last), ...volcano.slice(first, last), ...gold.slice(first, last), ...lime.slice(first, last), ...green.slice(first, last), ...cyan.slice(first, last), ...blue.slice(first, last), ...geekblue.slice(first, last), ...purple.slice(first, last), ...magenta.slice(first, last)];

    return colors[(Math.abs(i) + 1) % colors.length];
};