// Exports decompiled C for every function. Usage: -postScript ExportDecompiled.java <out.c>
// @category Export
import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.app.script.GhidraScript;
import ghidra.program.model.listing.Function;
import ghidra.program.model.listing.FunctionIterator;

import java.io.PrintWriter;

public class ExportDecompiled extends GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] args = getScriptArgs();
        String out = args.length > 0 ? args[0] : currentProgram.getName() + ".c";

        DecompInterface decomp = new DecompInterface();
        decomp.openProgram(currentProgram);

        int count = 0;
        try (PrintWriter w = new PrintWriter(out, "UTF-8")) {
            w.printf("// %s, image base %s%n%n", currentProgram.getName(), currentProgram.getImageBase());
            FunctionIterator it = currentProgram.getFunctionManager().getFunctions(true);
            while (it.hasNext() && !monitor.isCancelled()) {
                Function f = it.next();
                w.printf("// ==== %s @ %s (size %d)%n", f.getName(), f.getEntryPoint(), f.getBody().getNumAddresses());
                DecompileResults r = decomp.decompileFunction(f, 60, monitor);
                if (r != null && r.decompileCompleted()) {
                    w.println(r.getDecompiledFunction().getC());
                } else {
                    w.println("// decompilation failed");
                }
                count++;
            }
        }
        decomp.dispose();
        println("Exported " + count + " functions to " + out);
    }
}
