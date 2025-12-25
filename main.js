const a = 'f asdf';
const cmd = $`echo 1 "${a+a}" bbb`;
cmd.useStd(false, true, true);
cmd.run()
