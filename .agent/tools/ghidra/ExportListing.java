// Dumps disassembly for address ranges. Usage: -postScript ExportListing.java <out.txt> <start> <end> [<start> <end> ...]
// @category Export
import ghidra.app.script.GhidraScript;
import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Instruction;
import ghidra.program.model.listing.InstructionIterator;
import ghidra.program.model.symbol.Symbol;

import java.io.PrintWriter;

public class ExportListing extends GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] args = getScriptArgs();
        try (PrintWriter w = new PrintWriter(args[0], "UTF-8")) {
            for (int i = 1; i + 1 < args.length; i += 2) {
                Address start = toAddr(Long.decode(args[i]));
                Address end = toAddr(Long.decode(args[i + 1]));
                disassemble(start);
                w.printf("==== %s .. %s%n", start, end);
                InstructionIterator it = currentProgram.getListing().getInstructions(start, true);
                while (it.hasNext()) {
                    Instruction ins = it.next();
                    if (ins.getAddress().compareTo(end) > 0) {
                        break;
                    }
                    Symbol s = getSymbolAt(ins.getAddress());
                    if (s != null) {
                        w.printf("%s:%n", s.getName());
                    }
                    w.printf("  %s  %s%n", ins.getAddress(), ins);
                }
            }
        }
    }
}
