function check_input(content)  {
    if( content.trim() === "") return false;
    lines = content.split("\n");
    const regex = /^source\s*:\s*"[^"]+"\s*;\s*trigger\s*:\s*"[^"]*"\s*;\s*destination\s*:\s*"[^"]+"\s*;(?:\s*color\s*=\s*"[^"]+"\s*)?;?$/;
    for (let line of lines) {
        line = line.trim(); 
        if( line === "") continue;
        if( regex.test(line) === false ){
            console.log("not match\n");
        }
    }
    return true;
}