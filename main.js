const a = 1;
const cmd = $`echo 1 ${a+a} bbb`;
cmd.useStd(false, true, true);
cmd.run()
