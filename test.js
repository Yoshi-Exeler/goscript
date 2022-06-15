function main() {
    let a = getLoopAfter();
    console.log(a);
}

function getLoopAfter() {
    let b = 0;
    for (let i = 0; i < 1000000000; i++) {
        b = i;
    }
    return b;
}

main();